package repository_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/hhiroshell/git-client-tui/src/models"
)

var _ = Describe("RepositoryService", func() {
	Describe("DiscoverRepository", func() {
		Context("when in a valid git repository", func() {
			It("should discover the repository successfully", func() {
				repo, err := repoService.DiscoverRepository()

				Expect(err).NotTo(HaveOccurred())
				Expect(repo).NotTo(BeNil())
				Expect(repo.WorkingDir).To(Equal(tempRepoPath))
				Expect(repo.Path).To(Equal(filepath.Join(tempRepoPath, ".git")))
				Expect(repo.CurrentBranch).NotTo(BeEmpty())
				Expect(repo.IsClean).To(BeTrue())
				Expect(repo.Remotes).To(HaveLen(1))
				Expect(repo.Remotes[0].Name).To(Equal("origin"))
				Expect(repo.Remotes[0].URL).To(Equal("https://github.com/test/test-repo.git"))
			})
		})

		Context("when not in a git repository", func() {
			BeforeEach(func() {
				// Create a temporary directory that's not a git repo
				nonGitDir, err := os.MkdirTemp("", "non-git-*")
				Expect(err).NotTo(HaveOccurred())

				err = os.Chdir(nonGitDir)
				Expect(err).NotTo(HaveOccurred())

				DeferCleanup(func() {
					os.Chdir(tempRepoPath)
					os.RemoveAll(nonGitDir)
				})
			})

			It("should return an appropriate error", func() {
				repo, err := repoService.DiscoverRepository()

				Expect(err).To(HaveOccurred())
				Expect(repo).To(BeNil())

				var serviceErr *models.ServiceError
				Expect(err).To(BeAssignableToTypeOf(serviceErr))
			})
		})
	})

	Describe("GetStatus", func() {
		Context("when repository is clean", func() {
			It("should return clean status", func() {
				// First discover the repository
				_, err := repoService.DiscoverRepository()
				Expect(err).NotTo(HaveOccurred())

				status, err := repoService.GetStatus()

				Expect(err).NotTo(HaveOccurred())
				Expect(status).NotTo(BeNil())
				Expect(status.IsClean).To(BeTrue())
				Expect(status.StagedChanges).To(HaveLen(0))
				Expect(status.UnstagedChanges).To(HaveLen(0))
				Expect(status.Branch).NotTo(BeEmpty())
			})
		})

		Context("when repository has unstaged changes", func() {
			BeforeEach(func() {
				// Create a new file
				testFile := filepath.Join(tempRepoPath, "test.txt")
				err := os.WriteFile(testFile, []byte("test content"), 0644)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should detect unstaged changes", func() {
				// First discover the repository
				_, err := repoService.DiscoverRepository()
				Expect(err).NotTo(HaveOccurred())

				status, err := repoService.GetStatus()

				Expect(err).NotTo(HaveOccurred())
				Expect(status).NotTo(BeNil())
				Expect(status.IsClean).To(BeFalse())
				Expect(status.UnstagedChanges).To(HaveLen(1))
				Expect(status.UnstagedChanges[0].FilePath).To(Equal("test.txt"))
				Expect(status.UnstagedChanges[0].ChangeType).To(Equal(models.Added))
				Expect(status.UnstagedChanges[0].IsStaged).To(BeFalse())
			})
		})

		Context("when repository has staged changes", func() {
			BeforeEach(func() {
				// Create and stage a new file
				testFile := filepath.Join(tempRepoPath, "staged.txt")
				err := os.WriteFile(testFile, []byte("staged content"), 0644)
				Expect(err).NotTo(HaveOccurred())

				worktree, err := testRepo.Worktree()
				Expect(err).NotTo(HaveOccurred())

				_, err = worktree.Add("staged.txt")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should detect staged changes", func() {
				// First discover the repository
				_, err := repoService.DiscoverRepository()
				Expect(err).NotTo(HaveOccurred())

				status, err := repoService.GetStatus()

				Expect(err).NotTo(HaveOccurred())
				Expect(status).NotTo(BeNil())
				Expect(status.IsClean).To(BeFalse())
				Expect(status.StagedChanges).To(HaveLen(1))
				Expect(status.StagedChanges[0].FilePath).To(Equal("staged.txt"))
				Expect(status.StagedChanges[0].ChangeType).To(Equal(models.Added))
				Expect(status.StagedChanges[0].IsStaged).To(BeTrue())
			})
		})

		Context("when repository has modified files", func() {
			BeforeEach(func() {
				// Modify the existing README.md file
				readmePath := filepath.Join(tempRepoPath, "README.md")
				err := os.WriteFile(readmePath, []byte("# Modified Test Repository\n\nThis file has been modified.\n"), 0644)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should detect modified files", func() {
				// First discover the repository
				_, err := repoService.DiscoverRepository()
				Expect(err).NotTo(HaveOccurred())

				status, err := repoService.GetStatus()

				Expect(err).NotTo(HaveOccurred())
				Expect(status).NotTo(BeNil())
				Expect(status.IsClean).To(BeFalse())
				Expect(status.UnstagedChanges).To(HaveLen(1))
				Expect(status.UnstagedChanges[0].FilePath).To(Equal("README.md"))
				Expect(status.UnstagedChanges[0].ChangeType).To(Equal(models.Modified))
				Expect(status.UnstagedChanges[0].IsStaged).To(BeFalse())
			})
		})
	})

	Describe("Integration tests", func() {
		It("should work with a complete workflow", func() {
			// Discover repository
			repo, err := repoService.DiscoverRepository()
			Expect(err).NotTo(HaveOccurred())
			Expect(repo).NotTo(BeNil())

			// Check initial status
			status, err := repoService.GetStatus()
			Expect(err).NotTo(HaveOccurred())
			Expect(status.IsClean).To(BeTrue())

			// Create a new file
			newFile := filepath.Join(tempRepoPath, "workflow.txt")
			err = os.WriteFile(newFile, []byte("workflow test"), 0644)
			Expect(err).NotTo(HaveOccurred())

			// Check status with new file
			status, err = repoService.GetStatus()
			Expect(err).NotTo(HaveOccurred())
			Expect(status.IsClean).To(BeFalse())
			Expect(status.UnstagedChanges).To(HaveLen(1))

			// Stage the file
			worktree, err := testRepo.Worktree()
			Expect(err).NotTo(HaveOccurred())
			_, err = worktree.Add("workflow.txt")
			Expect(err).NotTo(HaveOccurred())

			// Check status with staged file
			status, err = repoService.GetStatus()
			Expect(err).NotTo(HaveOccurred())
			Expect(status.IsClean).To(BeFalse())
			Expect(status.StagedChanges).To(HaveLen(1))
			Expect(status.UnstagedChanges).To(HaveLen(0))
		})
	})
})