{{/*
Check if a path should be excluded.
Returns true if the given path starts with any prefix in the list.
*/}}
{{- define "isExcluded" -}}
{{- $path := index . 0 -}}
{{- $list := index . 1 -}}
{{- $excluded := false -}}
{{- range $item := $list }}
  {{- if hasPrefix $path $item }}
    {{- $excluded = true -}}
  {{- end }}
{{- end }}
{{- $excluded -}}
{{- end }}
