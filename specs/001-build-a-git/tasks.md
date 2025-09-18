# Tasks: Git Client TUI Implementation

**Input**: Design documents from `/specs/001-build-a-git/`
**Prerequisites**: plan.md, research.md, data-model.md, contracts/

## Setup Tasks

- [ ] T001 Create project structure per implementation plan
- [ ] T002 Initialize Go project with Cobra CLI and required dependencies
- [ ] T003 [P] Configure linting and formatting tools
- [ ] T004 [P] Create basic CLI scaffolding with Cobra

## Test First Tasks (TDD)

**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Repository Management Tests
- [ ] T005 [P] Contract test for RepositoryService.DiscoverRepository in tests/contract/repository/discover_test.go
- [ ] T006 [P] Contract test for RepositoryService.GetStatus in tests/contract/repository/status_test.go

### Branch Management Tests
- [ ] T007 [P] Contract test for BranchService.ListBranches in tests/contract/branch/list_test.go
- [ ] T008 [P] Contract test for BranchService.CreateBranch in tests/contract/branch/create_test.go
- [ ] T009 [P] Contract test for BranchService.SwitchBranch in tests/contract/branch/switch_test.go
- [ ] T010 [P] Contract test for BranchService.MergeBranch in tests/contract/branch/merge_test.go

### Staging Tests
- [ ] T011 [P] Contract test for StagingService.StageFile in tests/contract/staging/file_test.go
- [ ] T012 [P] Contract test for StagingService.StageHunks in tests/contract/staging/hunk_test.go
- [ ] T013 [P] Contract test for StagingService.StageLines in tests/contract/staging/line_test.go

### Commit Tests
- [ ] T014 [P] Contract test for CommitService.GetCommitHistory in tests/contract/commit/history_test.go
- [ ] T015 [P] Contract test for CommitService.CreateCommit in tests/contract/commit/create_test.go
- [ ] T016 [P] Contract test for CommitService.SquashCommits in tests/contract/commit/squash_test.go

### Integration Tests
- [ ] T017 [P] Integration test for file-level staging workflow in tests/integration/staging_file_test.go
- [ ] T018 [P] Integration test for line-level staging workflow in tests/integration/staging_line_test.go
- [ ] T019 [P] Integration test for branch creation and switching in tests/integration/branch_test.go
- [ ] T020 [P] Integration test for commit squashing in tests/integration/commit_squash_test.go

## Core Implementation Tasks

### Model Implementation
- [ ] T021 [P] Implement Repository model in src/models/repository.go
- [ ] T022 [P] Implement Branch model in src/models/branch.go
- [ ] T023 [P] Implement Commit model in src/models/commit.go
- [ ] T024 [P] Implement Change model in src/models/change.go
- [ ] T025 [P] Implement Hunk model in src/models/hunk.go
- [ ] T026 [P] Implement Line model in src/models/line.go

### Service Implementation
- [ ] T027 Implement RepositoryService interface in src/services/repository/repository.go
- [ ] T028 Implement BranchService interface in src/services/branch/branch.go
- [ ] T029 Implement StagingService interface in src/services/staging/staging.go
- [ ] T030 Implement CommitService interface in src/services/commit/commit.go
- [ ] T031 [P] Implement diff parsing utility in src/lib/diff/parser.go
- [ ] T032 [P] Implement patch generation utility in src/lib/patch/generator.go

### TUI Implementation
- [ ] T033 Implement base TUI application with Bubble Tea in src/tui/app.go
- [ ] T034 Implement file view component in src/tui/views/file_view.go
- [ ] T035 Implement diff view component in src/tui/views/diff_view.go
- [ ] T036 Implement branch view component in src/tui/views/branch_view.go
- [ ] T037 Implement commit view component in src/tui/views/commit_view.go
- [ ] T038 Implement status bar component in src/tui/components/status_bar.go
- [ ] T039 Implement keyboard input handling in src/tui/input.go
- [ ] T040 Implement theme and styling with Lip Gloss in src/tui/style/theme.go

## Integration Tasks

- [ ] T041 Integrate repository service with TUI in src/tui/app.go
- [ ] T042 Integrate branch service with branch view in src/tui/views/branch_view.go
- [ ] T043 Integrate staging service with file and diff views
- [ ] T044 Integrate commit service with commit view
- [ ] T045 Implement error handling and user notifications
- [ ] T046 Add context-sensitive help system

## Polish Tasks

- [ ] T047 [P] Add unit tests for diff parser in tests/unit/diff/parser_test.go
- [ ] T048 [P] Add unit tests for patch generator in tests/unit/patch/generator_test.go
- [ ] T049 [P] Performance optimization for large repositories
- [ ] T050 [P] Add terminal size detection and responsive layout
- [ ] T051 [P] Create user documentation in docs/usage.md
- [ ] T052 Implement debug logging system
- [ ] T053 Run all tests and fix any remaining issues
- [ ] T054 Run quickstart scenarios to verify functionality

## Dependencies

- Tests (T005-T020) before implementation (T021-T046)
- Model implementations (T021-T026) before corresponding services (T027-T030)
- Service implementations (T027-T030) before TUI integration (T041-T044)
- Core functionality before polish tasks

## Parallel Execution Examples

```
# Launch model implementation tasks in parallel:
Task: "Implement Repository model in src/models/repository.go"
Task: "Implement Branch model in src/models/branch.go"
Task: "Implement Commit model in src/models/commit.go"
Task: "Implement Change model in src/models/change.go"
Task: "Implement Hunk model in src/models/hunk.go"
Task: "Implement Line model in src/models/line.go"

# Launch contract tests in parallel:
Task: "Contract test for RepositoryService.DiscoverRepository in tests/contract/repository/discover_test.go"
Task: "Contract test for BranchService.ListBranches in tests/contract/branch/list_test.go"
Task: "Contract test for StagingService.StageFile in tests/contract/staging/file_test.go"
Task: "Contract test for CommitService.GetCommitHistory in tests/contract/commit/history_test.go"
```

## Notes

- [P] tasks can run in parallel (different files, no dependencies)
- Verify tests fail before implementing the corresponding functionality
- Follow the test-first development approach
- Ensure proper error handling in all service implementations
- Focus on performance optimization for large repositories