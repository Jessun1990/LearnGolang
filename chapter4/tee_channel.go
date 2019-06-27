chapter4

func teeChannelExample(){
	tee := func(
		done<-chan interface{}, in <- chan interface{})(
			_, _ <-chan interface{}){

		out1 := make(chan interface{})
		out2 := make(chan interface{})

		go func(){
			defer close(out1)
			defer close(out2)

		}
	}()
}
