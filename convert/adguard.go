package convert

func ConvertFromAdguard(sourceFile string, targetFile string) error {
	err := ConvertWithSingBox(sourceFile, targetFile, "adguard-srs")
	if err != nil {
		return err
	}
	return nil
}
