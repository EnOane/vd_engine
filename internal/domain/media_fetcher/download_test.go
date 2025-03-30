package media_fetcher

//func TestUploadToS3(t *testing.T) {
//	config.MustLoad()
//	s3.MustConnect()
//
//	in := make(chan []byte)
//
//	go func() {
//		defer close(in)
//
//		in <- []byte("text text/n")
//		in <- []byte("text text/n")
//		in <- []byte("text text/n")
//		in <- []byte("text text/n")
//	}()
//
//	err := uploadToS3(in, "text.txt")
//
//	assert.Nil(t, err)
//}
