{{/*
    Custom Reports ReactionListener CC v2
    
    Made By Devonte#0745 / Naru#6203
    Contributors: DZ#6669, Piter#5960, WickedWizard#3588
    
    Recommended Trigger Type: Reaction - Added Only
*/}}

{{/* THINGS TO CHANGE */}}

{{$staff := cslice ROLEID ROLEID}} {{/* A list of roles for people considered Admins. Replace ROLEID accordingly. */}}

{{$logChannel := }} {{/* Channel ID to log reports */}}

{{/* ACTUAL CODE - DO NOT TOUCH */}}

{{$isStaff := false}}
{{if .ReactionAdded}}
    {{if and .ReactionMessage.Author.Bot (eq .Reaction.ChannelID $logChannel)}}
        {{if or (dbGet 7 "reopen") (not .ReactionMessage.EditedTimestamp)}}
            {{range (getMember $.Reaction.UserID).Roles}}
                {{if in $staff .}}
                    {{$isStaff = true}}
                {{end}}
            {{end}}
            {{if $isStaff}}
                {{$mod := (userArg .Reaction.UserID).String}}
                {{$report := index (getMessage nil .ReactionMessage.ID).Embeds 0|structToSdict}}
                    {{range $k, $v := $report}}
                        {{if eq (kindOf $v true) "struct"}}
                            {{$report.Set $k (structToSdict $v)}}
                        {{end}}
                    {{end}}
                    {{$user := (userArg (reReplace `\(|\)` (reFind `\d{17,19}` $report.Description) "")).ID}}
                {{if eq .Reaction.Emoji.Name "✅"}}
                    {{with $report}}
                        {{.Set "color" 0x83fe25}}
                        {{.Set "description" (print "Report marked **Done** by " $mod " [\u200b](" $user ")")}}
                        {{.Set "timestamp" currentTime}}
                        {{.Author.Set "icon_url" .Author.IconURL}}
                    {{end}}
                    {{editMessage nil .ReactionMessage.ID (complexMessageEdit "embed" $report)}}
                    {{deleteAllMessageReactions nil .ReactionMessage.ID}}
                    {{dbDel 7 "reopen"}}
                {{else if eq .Reaction.Emoji.Name "❎"}}
                    {{with $report}}
                        {{.Set "color" 0xfeb225}}
                        {{.Set "description" (print "Report marked **Ignored** by " $mod " [\u200b](" $user ")")}}
                        {{.Set "timestamp" currentTime}}
                        {{.Author.Set "icon_url" .Author.IconURL}}
                    {{end}}
                    {{editMessage nil .ReactionMessage.ID (complexMessageEdit "embed" $report)}}
                    {{deleteAllMessageReactions nil .ReactionMessage.ID}}
                    {{dbDel 7 "reopen"}}
                {{else if eq .Reaction.Emoji.Name "🛡"}}
                    {{deleteAllMessageReactions nil .ReactionMessage.ID}}
                    {{addMessageReactions nil .ReactionMessage.ID "❌" "⚠" "🔇" "👢" "🔨"}}
                    {{dbSetExpire 7 "modaction" true 300}}
                {{else if eq .Reaction.Emoji.Name "❌" "⚠" "🔇" "👢" "🔨"}}
                    {{if (dbGet 7 "modaction")}}
                        {{$action := ""}}
                        {{if eq .Reaction.Emoji.Name "⚠"}}
                            {{if (reFind `(?i)ManageMessages` (exec "viewperms 204255221017214977"))}}
                                {{sendMessage nil "30s. Type out the reason why you want to warn that user."}}
				                {{sleep 30}}
				                {{$warn := (dbGet 0 "warn_kick_mute_ban").Value}}
                                {{$s := execAdmin "warn" $user $warn}}{{$action = "warned"}}
                            {{else}}
                                {{deleteMessageReaction nil .ReactionMessage.ID .Reaction.UserID "⚠"}}
                                {{print .User.Mention ", unable to warn the user: Missing Permissions `ManageMessages`"}}
                            {{end}}
                        {{else if eq .Reaction.Emoji.Name "🔇"}}
                            {{if (reFind `(?i)KickMembers` (exec "viewperms 204255221017214977"))}}
                                {{sendMessage nil "30s. Type out the reason why you want to mute that user."}}
				                {{sleep 30}}
				                {{$mute := (dbGet 0 "warn_kick_mute_ban").Value}}
                                {{$s := execAdmin "mute" $user $mute}}{{$action = "muted"}}
                            {{else}}
                                {{deleteMessageReaction nil .ReactionMessage.ID .Reaction.UserID "🔇"}}
                                {{print .User.Mention ", unable to mute the user: Missing Permissions `KickMembers`"}}
                            {{end}}
                        {{else if eq .Reaction.Emoji.Name "👢"}}
                            {{if (reFind `(?i)KickMembers` (exec "viewperms 204255221017214977"))}}
                                {{sendMessage nil "30s. Type out the reason why you want to kick that user."}}
				                {{sleep 30}}
				                {{$kick := (dbGet 0 "warn_kick_mute_ban").Value}}
                                {{$s := execAdmin "kick" $user $kick}}{{$action = "kicked"}}
                            {{else}}
                                {{deleteMessageReaction nil .ReactionMessage.ID .Reaction.UserID "👢"}}
                                {{print .User.Mention ", unable to kick the user: Missing Permissions `KickMembers`"}}
                            {{end}}
                        {{else if eq .Reaction.Emoji.Name "🔨"}}
                            {{if (reFind `(?i)BanMembers` (exec "viewperms 204255221017214977"))}}
                                {{sendMessage nil "30s. Type out the reason why you want to ban that user."}}
				                {{sleep 30}}
				                {{$ban := (dbGet 0 "warn_kick_mute_ban").Value}}
                                {{$s := execAdmin "BanMembers" $user $ban}}{{$action = "banned"}}
                            {{else}}
                                {{deleteMessageReaction nil .ReactionMessage.ID .Reaction.UserID "🔨"}}
                                {{print .User.Mention ", unable to ban the user: Missing Permissions `BanMembers`"}}
                            {{end}}
                        {{else if eq .Reaction.Emoji.Name "❌"}}
                            {{deleteAllMessageReactions nil .ReactionMessage.ID}}
                            {{addMessageReactions nil .ReactionMessage.ID "✅" "❎" "🛡"}}
                        {{end}}
                        {{if $action}}
                            {{with $report}}
                                {{.Set "color" 0xfe3025}}
                                {{.Set "description" (print "Mod-Action: **" $action "** by " $mod " [\u200b](" $user ")")}}
                                {{.Set "timestamp" currentTime}}
                                {{.Author.Set "icon_url" .Author.IconURL}}
                            {{end}}
                            {{editMessage nil .ReactionMessage.ID (complexMessageEdit "embed" $report)}}
                            {{deleteAllMessageReactions nil .ReactionMessage.ID}}
                            {{dbDel 7 "modaction"}}{{dbDel 7 "reopen"}}
                        {{end}}
                    {{else}}
                        {{deleteAllMessageReactions nil .ReactionMessage.ID}}
                        {{addMessageReactions nil .ReactionMessage.ID "✅" "❎" "🛡"}}
                    {{end}}
                {{end}}
            {{end}}
        {{end}}
    {{end}}
{{end}}
