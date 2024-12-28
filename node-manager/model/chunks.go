package model

type chunk struct {
	id     string
	fileID string
	nodeID string
	size   int64
}

// const (
// 	chunkSize = 2 * 1024 * 1024 // 2MB in bytes
// )

// file, handler, err := r.FormFile("myFile")
//         if err != nil {
//                 http.Error(w, "Unable to get file from form", http.StatusBadRequest)
//                 return
//         }
//         defer file.Close()

// func (dm *DbManager) CreateChunk(file io.Reader, handler *multipart.FileHeader) error {

// 	buf := make([]byte, chunkSize)
// 	chunkIndex := 0
// 	var chunks []chunk

// 	filename := handler.Filename
// 	ext := filepath.Ext(filename)
// 	baseName := strings.TrimSuffix(filename, ext)

// 	tx := dm.DB.Begin()

// 	object, err := dm.CreateFile(tx, filename, baseName)
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	for {
// 		bytesRead, err := file.Read(buf)
// 		if err != nil {
// 			if err != io.EOF {
// 				return fmt.Errorf("error reading file: %w", err)
// 			}
// 			break // Reached end of file
// 		}

// 		chunkName := fmt.Sprintf("%s.part%d%s", baseName, chunkIndex, ext)
// 		chunkData := make([]byte, bytesRead) // Create a new slice with the correct size
// 		copy(chunkData, buf[:bytesRead])     // Copy only the read bytes

// 		// node allocation and send data to node logic here
// 		chunks = append(chunks, chunk{
// 			id:     uuid.New().String(),
// 			fileID: object.id,
// 			nodeID: "nodeID",
// 			size:   int64(bytesRead),
// 		})

// 		fmt.Printf("Created chunk: %s, size: %d bytes\n", chunkName, bytesRead)
// 		chunkIndex++
// 	}
// 	return nil
// }

// func (dm *DbManager) CreateFile(tx *gorm.DB, filename string, baseName string) (file *object, err error) {

// 	object := &object{
// 		id:          uuid.New().String(),
// 		name:        filename,
// 		size:        int64(len(filename)),
// 		totalChunks: 0,
// 	}

// 	if err := tx.Create(object).Error; err != nil {
// 		tx.Rollback()
// 		return nil, err
// 	}

// 	return object, nil
// }
