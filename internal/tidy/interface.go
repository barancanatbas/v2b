package tidy

type TidyServiceInterface interface {
	UpdateModuleBranch(modulePath, branchName string) error
	UpdateModuleByID(moduleID int) error
}
