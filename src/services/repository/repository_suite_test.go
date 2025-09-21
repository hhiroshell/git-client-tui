package repository_test

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hhiroshell/git-client-tui/src/services/repository"
)

var (
	tempRepoPath   string
	testRepo       *git.Repository
	repoService    repository.RepositoryService
	originalWD     string
)

func TestRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Service Suite")
}

var _ = BeforeSuite(func() {
	var err error

	// Store the original working directory
	originalWD, err = os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	// Create a temporary directory for test repository
	tempRepoPath, err = os.MkdirTemp("", "git-client-tui-test-*")
	Expect(err).NotTo(HaveOccurred())

	// Initialize a git repository in the temporary directory
	testRepo, err = git.PlainInit(tempRepoPath, false)
	Expect(err).NotTo(HaveOccurred())

	// Configure the repository with test user
	cfg, err := testRepo.Config()
	Expect(err).NotTo(HaveOccurred())

	cfg.User.Name = "Test User"
	cfg.User.Email = "test@example.com"

	// Create an initial commit to have a proper HEAD
	worktree, err := testRepo.Worktree()
	Expect(err).NotTo(HaveOccurred())

	// Create a test file
	testFilePath := filepath.Join(tempRepoPath, "README.md")
	err = os.WriteFile(testFilePath, []byte("# Test Repository\n\nThis is a test repository for unit tests.\n"), 0644)
	Expect(err).NotTo(HaveOccurred())

	// Add the file to the index
	_, err = worktree.Add("README.md")
	Expect(err).NotTo(HaveOccurred())

	// Create the initial commit
	commit, err := worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
		},
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(commit).NotTo(Equal(plumbing.ZeroHash))

	// Add a remote for testing
	_, err = testRepo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/test/test-repo.git"},
	})
	Expect(err).NotTo(HaveOccurred())

	// Change to the test repository directory
	err = os.Chdir(tempRepoPath)
	Expect(err).NotTo(HaveOccurred())

	// Initialize the repository service
	repoService = repository.NewGitRepositoryService()
})

var _ = AfterSuite(func() {
	// Change back to the original working directory
	if originalWD != "" {
		err := os.Chdir(originalWD)
		Expect(err).NotTo(HaveOccurred())
	}

	// Clean up the temporary repository
	if tempRepoPath != "" {
		err := os.RemoveAll(tempRepoPath)
		Expect(err).NotTo(HaveOccurred())
	}
})

var _ = BeforeEach(func() {
	// Ensure we're in the test repository directory for each test
	err := os.Chdir(tempRepoPath)
	Expect(err).NotTo(HaveOccurred())

	// Reset the repository to a clean state
	worktree, err := testRepo.Worktree()
	Expect(err).NotTo(HaveOccurred())

	// Clean any uncommitted changes
	status, err := worktree.Status()
	Expect(err).NotTo(HaveOccurred())

	for filePath, fileStatus := range status {
		if fileStatus.Worktree == git.Modified || fileStatus.Worktree == git.Deleted {
			err = worktree.Checkout(&git.CheckoutOptions{
				Force: true,
			})
			Expect(err).NotTo(HaveOccurred())
			break
		}
		if fileStatus.Worktree == git.Untracked {
			err = os.Remove(filepath.Join(tempRepoPath, filePath))
			if err != nil && !os.IsNotExist(err) {
				Expect(err).NotTo(HaveOccurred())
			}
		}
	}
})

var _ = AfterEach(func() {
	// Additional cleanup after each test if needed
	// This ensures tests don't interfere with each other
})