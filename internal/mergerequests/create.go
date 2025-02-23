package mergerequests

import gitlab "gitlab.com/gitlab-org/api/client-go"

func createMergeRequestNote(repo string, mr int, comment string, client *gitlab.Client) error {
	noteOptions := &gitlab.CreateMergeRequestNoteOptions{
		Body: gitlab.Ptr(comment),
	}
	_, _, err := client.Notes.CreateMergeRequestNote(repo, mr, noteOptions)

	return err
}
