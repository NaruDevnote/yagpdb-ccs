{{/*
    Tags Main CC

    Made By Devonte#0745 / Naru#6203
    Contributors: Maverick Wolf#1010, TheHDCrafter#0001

    Recommended Trigger Type: Starts With
    Recommended Trigger     : ;

    Make sure that the trigger is a SEMI-COLON. See README.md for more info.

    © NaruDevnote 2020-2021 (GNU GPL v3)
    https://github.com/NaruDevnote/yagpdb-ccs
*/}}

{{/* THINGS TO CHANGE */}}

{{$staff := cslice ROLE-ID ROLE-ID}} {{/* Role IDs of staff roles (at least one is required). */}}

{{/* ACTUAL CODE - DO NOT TOUCH */}}

{{if .CmdArgs}}
    {{if gt (len .CmdArgs) 1}}
        {{$t := (index .CmdArgs 0|lower)}}{{$sub := (index .CmdArgs 1|lower)}}
        {{if and (eq $t "tag") (eq $sub "add" "edit" "del" "delete")}}
            {{$isStaff := false}}
            {{range $staff}}{{if hasRoleID .}}{{$isStaff = true}}{{end}}{{end}}
            {{if $isStaff}}
                {{if (eq $sub "add")}}
                    {{$a := parseArgs 5 "Usage: `;tag add <name/aliases:text> [color:decimal] [image:URL] <content>`" (carg "string" "") (carg "string" "") (carg "string" "name") (carg "string" "color") (carg "string" "image") (carg "string" "content")}}
                    {{if (reFind `(?i)[^a-zA-Z0-9/-]+|tag` ($a.Get 2))}}You cant have numbers or special characters in the tag name.{{else}}
                    {{if (dbGet 0 (print "tag_" (lower ($a.Get 2))))}}A tag with this name already exists! Use `;tag edit` instead{{else}}
                    {{dbSet 0 (print "tag_" (lower ($a.Get 2))) (sdict "author" .User.String "color" (or (toInt ($a.Get 3)) 0) "content" (or ($a.Get 5) "") "image" (or ($a.Get 4) ""))}}Added `{{lower ($a.Get 2)}}`{{end}}{{end}}
                {{else if (eq $sub "edit")}}
                    {{$a := parseArgs 3 "Usage: `;tag edit <full-name> [new-aliases:text] [color:decimal] [image:URL] [new-content]`" (carg "string" "") (carg "string" "") (carg "string" "name") (carg "string" "aliases") (carg "string" "color") (carg "string" "image") (carg "string" "content")}}
                    {{with (dbGet 0 (print "tag_" (lower ($a.Get 2))))}}
                        {{$nt := sdict .Value}}{{$nn := .Key}}
                        {{if ($a.IsSet 4)}}{{if (reFind `\A[0-9]+\z` ($a.Get 4))}}{{$nt.Set "color" (toInt ($a.Get 4))}}{{end}}{{end}}
                        {{if ($a.IsSet 5)}}{{$nt.Set "image" ($a.Get 5)}}{{end}}
                        {{if ($a.IsSet 6)}}{{$nt.Set "content" ($a.Get 6)}}{{end}}
                        {{if ($a.IsSet 3)}}
                            {{if (reFind `(?i)[^a-zA-Z0-9/-]+|tag` ($a.Get 3))}}You cant have numbers or special characters as aliases.{{else}}
                            {{$nn = print "snippet_" (lower ($a.Get 3))}}{{$nn = reReplace `/\z` $nn ""}}
							{{$nt.Set "author" $.User.String}}{{dbDel .UserID .Key}}
							{{dbSet 0 $nn $nt}}Edited `{{lower ($a.Get 2)}}`{{end}}
                        {{end}}
                    {{else}}Unknown tag (specify the whole name).{{end}}
                {{else if (eq $sub "del" "delete")}}
                    {{$a := parseArgs 3 "Usage: `;tag delete <name>`\nThis has to be the full tag name." (carg "string" "") (carg "string" "") (carg "string" "name")}}
                    {{with (dbGet 0 (print "tag_" (lower ($a.Get 2))))}}
                        {{dbDel .UserID .Key}}Deleted `{{lower ($a.Get 2)}}`
                    {{else}}That tag doesn't exist.{{end}}
                {{end}}
            {{end}}
        {{else if and (eq $t "tag") (eq $sub "list" "search" "dump" "help")}}
            {{if eq $sub "list"}}{{$list := ""}}
                {{range (dbTopEntries `tag\_%` 10 0)}}{{$list = print $list "\n`" (slice .Key 4) "`"}}{{else}}{{$list = "No Tags"}}{{end}}
                {{$msg := sendMessageRetID nil (cembed "title" "Tag List" "description" $list "footer" (sdict "text" "React with 🗑 to delete this message.\u200b\nPage: 1"))}}
                {{addMessageReactions nil $msg "🗑" "◀" "▶"}}
            {{else if eq $sub "search"}}
                {{$a := parseArgs 3 "Usage: `;tag search <name>`" (carg "string" "") (carg "string" "") (carg "string" "name")}}
                {{$res := ""}}
                {{range (dbTopEntries (print `tag\_%` (lower ($a.Get 2)) `%`) 50 0)}}{{$res = print $res "\n`" (slice .Key 4) "`"}}{{else}}{{$res = "No Results"}}{{end}}
                {{$msg := sendMessageRetID nil (cembed "title" "Tag Search Results" "description" $res "footer" (sdict "text" "React with 🗑 to delete this message.\u200b"))}}
                {{addMessageReactions nil $msg "🗑"}}
            {{else if eq $sub "dump"}}
                {{$a := parseArgs 3 "Usage: `;tag dump <name>`\nThis has to be the full tag name." (carg "string" "") (carg "string" "") (carg "string" "name")}}
                {{with (dbGet 0 (print "tag_" (lower ($a.Get 2))))}}
                    {{$dump := printf "GUILD: %s (%d)\nTAG NAME: %s\nREQUESTED BY: %s (%d)\n\n%s" $.Guild.Name $.Guild.ID .Key $.User.Username $.User.ID (reReplace `,` (json .Value) ",\n")}}
                    {{sendMessage nil (complexMessage "content" (print "Tag Info Dump for: **" (slice .Key 4|title) "**") "file" $dump)}}
                {{else}}Unknown tag (specify the whole name).{{end}}
            {{else if eq $sub "help"}}
                {{$embed := cembed
                    "title" "Tag Help"
                    "description" "**Key:**\n**<>** Required Args - **[]** Optional Args\nFor default values, put empty quotes \"\" and/or 0 for color.\n\n`;tagname` - Sends the tag\n`;tag add ...` - Adds a tag to the database (under the key `tag_` + tag-name).\n`;tag edit ...` - Edits an existing tag with the specified fields.\n`;tag dump <name>` - Sends the JSON content of a tag\n`;tag del/delete <name>` - Deletes an existing tag. This has to be the full tag name, not an alias.\n`;tag search <name>` - Searches for a tag based on the name/alias provided.\n`;tag list` - Lists all the tags in the server.\n`;tag dump <name>` - Sends a .txt attachment of the raw tag`\n`;tag help` - Sends this message. :)"
                    "footer" (sdict "text" "React with 🗑 to delete this message.\u200b")
                    "timestamp" currentTime}}
                {{$msg := sendMessageRetID nil $embed}}{{addMessageReactions nil $msg "🗑"}}
            {{end}}
        {{end}}
    {{else}}
        {{$tag := ""}}
        {{range (dbTopEntries (print `tag\_%` (lower (index .CmdArgs 0)) `%`) 100 0)}}{{$tag = print $tag (slice .Key 4)}}{{end}}
        {{with (dbGet 0 (print "tag_" $tag))}}{{$i := sdict .Value}}
            {{$msg := sendMessageRetID nil (cembed "title" (print "Tag: " (slice .Key 4|title)) "description" $i.content "color" (or $i.color 0) "image" (sdict "url" (or $i.image "")) "footer" (sdict "text" (print "Author: " (or $i.author "Could Not Find") " (from last edit)\nReact with 📱 to be DMed a mobile version.\u200b")))}}
            {{addMessageReactions nil $msg "📱"}}
        {{end}}
    {{end}}
{{end}}
