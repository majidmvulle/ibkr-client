{{/*
Expand the name of the chart.
*/}}
{{- define "ibkr-client.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "ibkr-client.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
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
{{- define "ibkr-client.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "ibkr-client.labels" -}}
helm.sh/chart: {{ include "ibkr-client.chart" . }}
{{ include "ibkr-client.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "ibkr-client.selectorLabels" -}}
app.kubernetes.io/name: {{ include "ibkr-client.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "ibkr-client.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "ibkr-client.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Image name for main application
*/}}
{{- define "ibkr-client.image" -}}
{{- $tag := .Values.image.tag | default .Chart.AppVersion }}
{{- printf "%s:%s" .Values.image.repository $tag }}
{{- end }}

{{/*
Image name for migrations
*/}}
{{- define "ibkr-client.migrationsImage" -}}
{{- $tag := .Values.migrations.image.tag | default .Chart.AppVersion }}
{{- printf "%s:%s" .Values.migrations.image.repository $tag }}
{{- end }}

{{/*
Image name for IBKR Gateway
*/}}
{{- define "ibkr-client.ibkrGatewayImage" -}}
{{- $tag := .Values.ibkrGateway.image.tag | default .Chart.AppVersion }}
{{- printf "%s:%s" .Values.ibkrGateway.image.repository $tag }}
{{- end }}
