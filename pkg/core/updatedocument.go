package core

// UpdateDocument inserts or updates a document in the zinc index
func (index *Index) UpdateDocument(docID string, doc *map[string]interface{}, mintedID bool) error {

	bdoc, err := index.BuildBlugeDocumentFromJSON(docID, doc)

	if err != nil {
		return err
	}

	// Finally update the document on disk
	writer := index.Writer
	if !mintedID {
		err = writer.Update(bdoc.ID(), bdoc)
	} else {
		err = writer.Insert(bdoc)
		index.GainDocsCount(1)
	}
	// err = writer.Update(bdoc.ID(), bdoc)
	if err != nil {
		return err
	}

	return nil

}
