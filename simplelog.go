package alog

// var _file PrintLogger

// func getFile() (PrintLogger, error) {
// 	if _file == nil {
// 		filename := filepath.Base(os.Args[0]) + ".log"
// 		directory := filepath.Join(os.Getenv("HOME"), "log")
// 		path := filepath.Join(directory, filename)
// 		var err error
// 		_file, err = NewRotatingLogger(path, DefaultLogger)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return _file, nil
// }

// func Log(msg string) error {
// 	file, err := getFile()
// 	if err != nil {
// 		return err
// 	}
// 	file.Println(strings.TrimRight(msg, "\n"))
// 	return nil
// }

// func Logf(msg string, args ...interface{}) error {
// 	return Log(fmt.Sprintf(Colorify(msg), args...))
// }
