package core

func (index *Index) UpdateDocument(docID string, doc *map[string]interface{}) error {

	bdoc, err := index.GetBlugeDocument(docID, doc)

	if err != nil {
		return err
	}

	// Finally update the document on disk
	writer := index.Writer
	err = writer.Update(bdoc.ID(), bdoc)
	if err != nil {
		return err
	}

	return nil

}
