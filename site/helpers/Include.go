package helpers

import "path/filepath"

func Include(path string) []string {
	files, _ := filepath.Glob("site/views/templates/*.html")

	// Ana klasördeki dosyaları ara
	path_files, _ := filepath.Glob("site/views/" + path + "/*.html")

	// Alt klasörlerdeki dosyaları da ara
	subdir_files, _ := filepath.Glob("site/views/" + path + "/*/*.html")

	// Tüm dosyaları birleştir
	files = append(files, path_files...)
	files = append(files, subdir_files...)

	return files
}
