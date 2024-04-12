package engine

type Camera interface {
	GetProjectionMat()
	GetViewMat()
}
