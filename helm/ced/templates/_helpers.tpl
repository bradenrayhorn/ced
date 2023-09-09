{{/*
Expand the name of the chart.
*/}}
{{- define "ced.name" -}}
{{- default "ced" .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "ced.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default "ced" .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "ced.chart" -}}
{{- printf "%s-%s" "ced" .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "ced.labels" -}}
helm.sh/chart: {{ include "ced.chart" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "ced.labelsServer" -}}
{{ include "ced.selectorLabelsServer" . }}
{{ include "ced.labels" . }}
{{- end }}

{{- define "ced.labelsUI" -}}
{{ include "ced.selectorLabelsUI" . }}
{{ include "ced.labels" . }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "ced.selectorLabelsServer" -}}
app.kubernetes.io/name: {{ include "ced.name" . }}-server
app.kubernetes.io/part-of: {{ include "ced.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
{{- define "ced.selectorLabelsUI" -}}
app.kubernetes.io/name: {{ include "ced.name" . }}-ui
app.kubernetes.io/part-of: {{ include "ced.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "ced.serviceAccountName" -}}
{{- default (include "ced.fullname" .) .Values.serviceAccount.name }}
{{- end }}
