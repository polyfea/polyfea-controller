package polyfea

func strToPtr(in string) *string {
	if in != "" {
		return &in
	}
	return nil
}

func arrToPtr[T any](in []T) *[]T {
	if in != nil {
		return &in
	}
	return nil
}
