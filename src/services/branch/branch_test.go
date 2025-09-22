package branch_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/hhiroshell/git-client-tui/src/models"
	"github.com/hhiroshell/git-client-tui/src/services/branch"
)

var _ = Describe("BranchService", func() {
	var (
		branchService branch.BranchService
	)

	BeforeEach(func() {
		// This will fail until we implement the BranchService
		branchService = branch.NewGitBranchService(testRepo)
	})

	Describe("ListBranches", func() {
		Context("when repository has multiple branches", func() {
			BeforeEach(func() {
				// Create additional test branches
				worktree, err := testRepo.Worktree()
				Expect(err).NotTo(HaveOccurred())

				// Create feature branch
				err = testRepo.CreateBranch(&config.Branch{
					Name:   "feature-branch",
					Remote: "",
					Merge:  plumbing.ReferenceName("refs/heads/feature-branch"),
				})
				Expect(err).NotTo(HaveOccurred())

				// Create a commit on feature branch
				err = worktree.Checkout(&git.CheckoutOptions{
					Branch: plumbing.ReferenceName("refs/heads/feature-branch"),
					Create: true,
				})
				Expect(err).NotTo(HaveOccurred())

				// Create a test file on feature branch
				featureFilePath := filepath.Join(tempRepoPath, "feature.txt")
				err = os.WriteFile(featureFilePath, []byte("Feature content"), 0644)
				Expect(err).NotTo(HaveOccurred())

				_, err = worktree.Add("feature.txt")
				Expect(err).NotTo(HaveOccurred())

				_, err = worktree.Commit("Add feature", &git.CommitOptions{
					Author: &object.Signature{
						Name:  "Test User",
						Email: "test@example.com",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				// Switch back to main branch
				err = worktree.Checkout(&git.CheckoutOptions{
					Branch: plumbing.ReferenceName("refs/heads/main"),
				})
				if err != nil {
					// Try master if main doesn't exist
					err = worktree.Checkout(&git.CheckoutOptions{
						Branch: plumbing.ReferenceName("refs/heads/master"),
					})
				}
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return success with branch list", func() {
				branchList, err := branchService.ListBranches()

				Expect(err).NotTo(HaveOccurred())
				Expect(branchList).NotTo(BeNil())
				Expect(branchList.Current).To(Or(Equal("main"), Equal("master")))
				Expect(len(branchList.Local)).To(BeNumerically(">=", 2))

				// Check that we have main/master and feature-branch
				branchNames := make([]string, len(branchList.Local))
				for i, branch := range branchList.Local {
					branchNames[i] = branch.Name
				}
				Expect(branchNames).To(ContainElements("feature-branch"))
			})
		})

		Context("when git command fails", func() {
			var (
				corruptedRepo *git.Repository
			)

			BeforeEach(func() {
				// Create a corrupted repository scenario
				corruptedRepoPath, err := os.MkdirTemp("", "corrupted-repo-*")
				Expect(err).NotTo(HaveOccurred())

				// Create a directory that looks like a git repo but is corrupted
				gitDir := filepath.Join(corruptedRepoPath, ".git")
				err = os.MkdirAll(gitDir, 0755)
				Expect(err).NotTo(HaveOccurred())

				// Write invalid git data
				err = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("invalid-ref"), 0644)
				Expect(err).NotTo(HaveOccurred())

				// Try to open it (this will likely fail)
				corruptedRepo, err = git.PlainOpen(corruptedRepoPath)
				if err == nil {
					branchService = branch.NewGitBranchService(corruptedRepo)
				}

				DeferCleanup(func() {
					os.RemoveAll(corruptedRepoPath)
				})
			})

			It("should return ErrGitCommandFailed", func() {
				if corruptedRepo == nil {
					Skip("Could not create corrupted repo scenario")
				}

				branchList, err := branchService.ListBranches()

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring("git command failed")))
				Expect(branchList).To(BeNil())
			})
		})
	})

	Describe("CreateBranch", func() {
		Context("when creating a new branch with valid name", func() {
			It("should create branch successfully", func() {
				err := branchService.CreateBranch("new-feature")

				Expect(err).NotTo(HaveOccurred())

				// Verify branch was created
				branchList, err := branchService.ListBranches()
				Expect(err).NotTo(HaveOccurred())

				branchNames := make([]string, len(branchList.Local))
				for i, branch := range branchList.Local {
					branchNames[i] = branch.Name
				}
				Expect(branchNames).To(ContainElement("new-feature"))
			})
		})

		Context("when creating a branch with invalid name", func() {
			It("should return ErrInvalidBranchName", func() {
				err := branchService.CreateBranch("invalid..name")

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(models.ErrInvalidBranchName))
			})
		})

		Context("when creating a branch that already exists", func() {
			BeforeEach(func() {
				// Create a branch first
				err := branchService.CreateBranch("existing-branch")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return ErrBranchExists", func() {
				err := branchService.CreateBranch("existing-branch")

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(models.ErrBranchExists))
			})
		})
	})

	Describe("SwitchBranch", func() {
		Context("when switching to an existing branch", func() {
			BeforeEach(func() {
				// Create a test branch first
				err := branchService.CreateBranch("test-branch")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should switch branch successfully", func() {
				err := branchService.SwitchBranch("test-branch")

				Expect(err).NotTo(HaveOccurred())

				// Verify current branch changed
				branchList, err := branchService.ListBranches()
				Expect(err).NotTo(HaveOccurred())
				Expect(branchList.Current).To(Equal("test-branch"))
			})
		})

		Context("when switching to a non-existent branch", func() {
			It("should return ErrBranchNotFound", func() {
				err := branchService.SwitchBranch("non-existent-branch")

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(models.ErrBranchNotFound))
			})
		})

		Context("when there are uncommitted changes", func() {
			BeforeEach(func() {
				// Create uncommitted changes
				testFilePath := filepath.Join(tempRepoPath, "uncommitted.txt")
				err := os.WriteFile(testFilePath, []byte("Uncommitted content"), 0644)
				Expect(err).NotTo(HaveOccurred())

				// Create a target branch
				err = branchService.CreateBranch("target-branch")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return ErrUncommittedChanges", func() {
				err := branchService.SwitchBranch("target-branch")

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(models.ErrUncommittedChanges))
			})
		})
	})

	Describe("MergeBranch", func() {
		Context("when merging a branch successfully", func() {
			BeforeEach(func() {
				// Create and switch to feature branch
				err := branchService.CreateBranch("merge-feature")
				Expect(err).NotTo(HaveOccurred())

				err = branchService.SwitchBranch("merge-feature")
				Expect(err).NotTo(HaveOccurred())

				// Create a commit on feature branch
				worktree, err := testRepo.Worktree()
				Expect(err).NotTo(HaveOccurred())

				mergeFilePath := filepath.Join(tempRepoPath, "merge-feature.txt")
				err = os.WriteFile(mergeFilePath, []byte("Merge feature content"), 0644)
				Expect(err).NotTo(HaveOccurred())

				_, err = worktree.Add("merge-feature.txt")
				Expect(err).NotTo(HaveOccurred())

				_, err = worktree.Commit("Add merge feature", &git.CommitOptions{
					Author: &object.Signature{
						Name:  "Test User",
						Email: "test@example.com",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				// Switch back to main branch
				err = branchService.SwitchBranch("main")
				if err != nil {
					err = branchService.SwitchBranch("master")
				}
				Expect(err).NotTo(HaveOccurred())
			})

			It("should merge branch successfully", func() {
				mergeResult, err := branchService.MergeBranch("merge-feature")

				Expect(err).NotTo(HaveOccurred())
				Expect(mergeResult).NotTo(BeNil())
				Expect(mergeResult.CommitHash).NotTo(BeEmpty())
				Expect(mergeResult.FilesChanged).To(ContainElement("merge-feature.txt"))
				Expect(mergeResult.HasConflicts()).To(BeFalse())
			})
		})

		Context("when fast-forward merge is possible", func() {
			BeforeEach(func() {
				// Create and switch to feature branch
				err := branchService.CreateBranch("fast-forward-feature")
				Expect(err).NotTo(HaveOccurred())

				err = branchService.SwitchBranch("fast-forward-feature")
				Expect(err).NotTo(HaveOccurred())

				// Create a commit on feature branch
				worktree, err := testRepo.Worktree()
				Expect(err).NotTo(HaveOccurred())

				ffFilePath := filepath.Join(tempRepoPath, "fast-forward.txt")
				err = os.WriteFile(ffFilePath, []byte("Fast forward content"), 0644)
				Expect(err).NotTo(HaveOccurred())

				_, err = worktree.Add("fast-forward.txt")
				Expect(err).NotTo(HaveOccurred())

				_, err = worktree.Commit("Add fast forward feature", &git.CommitOptions{
					Author: &object.Signature{
						Name:  "Test User",
						Email: "test@example.com",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				// Switch back to main branch
				err = branchService.SwitchBranch("main")
				if err != nil {
					err = branchService.SwitchBranch("master")
				}
				Expect(err).NotTo(HaveOccurred())
			})

			It("should perform fast-forward merge", func() {
				mergeResult, err := branchService.MergeBranch("fast-forward-feature")

				Expect(err).NotTo(HaveOccurred())
				Expect(mergeResult).NotTo(BeNil())
				Expect(mergeResult.FastForward).To(BeTrue())
				Expect(mergeResult.FilesChanged).To(ContainElement("fast-forward.txt"))
			})
		})

		Context("when merge has conflicts", func() {
			BeforeEach(func() {
				// Create conflicting changes on main branch
				worktree, err := testRepo.Worktree()
				Expect(err).NotTo(HaveOccurred())

				conflictFilePath := filepath.Join(tempRepoPath, "conflict.txt")
				err = os.WriteFile(conflictFilePath, []byte("Main branch content"), 0644)
				Expect(err).NotTo(HaveOccurred())

				_, err = worktree.Add("conflict.txt")
				Expect(err).NotTo(HaveOccurred())

				_, err = worktree.Commit("Add conflict file on main", &git.CommitOptions{
					Author: &object.Signature{
						Name:  "Test User",
						Email: "test@example.com",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				// Create and switch to feature branch
				err = branchService.CreateBranch("conflict-feature")
				Expect(err).NotTo(HaveOccurred())

				err = branchService.SwitchBranch("conflict-feature")
				Expect(err).NotTo(HaveOccurred())

				// Create conflicting change on feature branch
				err = os.WriteFile(conflictFilePath, []byte("Feature branch content"), 0644)
				Expect(err).NotTo(HaveOccurred())

				_, err = worktree.Add("conflict.txt")
				Expect(err).NotTo(HaveOccurred())

				_, err = worktree.Commit("Modify conflict file on feature", &git.CommitOptions{
					Author: &object.Signature{
						Name:  "Test User",
						Email: "test@example.com",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				// Switch back to main branch
				err = branchService.SwitchBranch("main")
				if err != nil {
					err = branchService.SwitchBranch("master")
				}
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return ErrMergeConflict", func() {
				mergeResult, err := branchService.MergeBranch("conflict-feature")

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(models.ErrMergeConflict))
				Expect(mergeResult).NotTo(BeNil())
				Expect(mergeResult.HasConflicts()).To(BeTrue())
				Expect(mergeResult.Conflicts).To(ContainElement("conflict.txt"))
			})
		})
	})

	Describe("FetchRemotes", func() {
		Context("when remote is reachable", func() {
			It("should fetch successfully", func() {
				Skip("Skipping remote fetch test - requires network access")

				err := branchService.FetchRemotes()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when remote is unreachable", func() {
			It("should return ErrFetchFailed", func() {
				// This test would require mocking network failure
				// For now, we'll test the interface exists
				err := branchService.FetchRemotes()

				// Since we're using a fake remote URL, this should fail
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring("fetch operation failed")))
			})
		})
	})

	Describe("Error Handling", func() {
		Context("when branch service operations fail", func() {
			It("should wrap errors with proper context", func() {
				// Test error wrapping by calling operations that should fail
				err := branchService.SwitchBranch("non-existent")

				Expect(err).To(HaveOccurred())
				// Verify error can be identified through errors.Is
				Expect(err).To(MatchError(models.ErrBranchNotFound))
			})
		})
	})
})