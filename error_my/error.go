package error_my
import (
	"errors"
)
var (
	UserConflict = errors.New("Người dùng đã tồn tại")
	SignUpFail = errors.New("Đăng kí người dùng thất bại")
	UserNotFound =errors.New("Không tìm thấy người dùng/Không tồn tại")
	UserUpdateFail =errors.New("Cập nhật thông tin người dùng thất bại")

	RepoConflict = errors.New("RepoGit này đã tồn tại")
	RepoNotFound =errors.New("Không tìm thấy RepoGit/Không tồn tại")
	RepoUpdateFail =errors.New("Cập nhật thông tin RepoGit thất bại")

	BookMarkConflict = errors.New("Bookmark này đã tồn tại")
	BookMarkNotFound =errors.New("Không tìm thấy Bookmark/Không tồn tại")
	DeleteBookMarkFail =errors.New("Xóa Bookmark thất bại")
)