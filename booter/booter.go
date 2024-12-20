package booter

type CleanUpFunc func()

type BootFunc func(c *Container) CleanUpFunc

func Boot(boots []BootFunc) CleanUpFunc {
	c := NewContainer()
	ch := make(chan CleanUpFunc, len(boots))

	for _, boot := range boots {
		go func() {
			ch <- boot(c)
		}()
	}

	return func() {
		for i := 0; i < len(boots); i++ {
			clean := <-ch
			if clean != nil {
				clean()
			}
		}
	}
}
