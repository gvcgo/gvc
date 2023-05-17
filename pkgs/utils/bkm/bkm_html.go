package bkm

const (
	Dt string = `<DT>%s`

	A string = `<A HREF="%s" ADD_DATE="%s" ICON="">%s</A>`

	// <h3 add_date="1675318038" last_modified="1679754685" personal_toolbar_folder="true">书签栏</h3>
	H3      string = `<H3 ADD_DATE="%s" LAST_MODIFIED="%s"%s>%s</H3>`
	TOOLBAR string = `PERSONAL_TOOLBAR_FOLDER="true"`

	P string = `<p>`

	Header string = `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
		It will be read and overwritten.
		DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>`

	Folder string = `<DT>%s
<DL><p>
%s
</DL><p>`
)
