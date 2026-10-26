import { selectedFiles } from './api'
import { useId, useRef, useState } from 'react'
import type { DocumentRequirement } from './creditData'
import { fileValidationError } from './creditData'
import { Icon } from './ui'

export function DocumentUploads({
  requirements,
  files,
  onChange,
  attempted = false,
}: {
  requirements: DocumentRequirement[]
  files: Record<string, string>
  onChange: (files: Record<string, string>) => void
  attempted?: boolean
}) {
  return (
    <div className="document-uploads">
      {requirements.map((requirement) => (
        <DocumentUpload
          key={requirement.id}
          requirement={requirement}
          filename={files[requirement.id]}
          attempted={attempted}
          onChange={(filename) => {
            const next = { ...files }
            if (filename) next[requirement.id] = filename
            else delete next[requirement.id]
            onChange(next)
          }}
        />
      ))}
    </div>
  )
}

function DocumentUpload({
  requirement,
  filename,
  onChange,
  attempted,
}: {
  requirement: DocumentRequirement
  filename?: string
  onChange: (name?: string) => void
  attempted: boolean
}) {
  const input = useRef<HTMLInputElement>(null)
  const [error, setError] = useState('')
  const id = useId()
  const missing = attempted && requirement.required && !filename
  return (
    <div
      className={`document-upload ${filename ? 'is-uploaded' : ''} ${missing || error ? 'has-error' : ''}`}
    >
      <input
        ref={input}
        id={id}
        className="sr-only"
        type="file"
        tabIndex={-1}
        aria-label={`Файл: ${requirement.title}`}
        accept=".pdf,.png,.jpg,.jpeg"
        onChange={(event) => {
          const file = event.target.files?.[0]
          event.target.value = ''
          if (!file) return
          const message = fileValidationError(file)
          setError(message)
          if (!message) {
            selectedFiles.set(requirement.id, file)
            onChange(file.name)
          }
        }}
      />
      <button
        type="button"
        className="document-upload-button"
        aria-label={`Прикрепить: ${requirement.title}${requirement.required ? ' *' : ''} ${filename || 'Нажмите, чтобы прикрепить документ'}`}
        aria-describedby={`${id}-description`}
        aria-invalid={missing || !!error || undefined}
        onClick={() => input.current?.click()}
      >
        <span className="document-icon">
          <Icon name="doc" />
        </span>
        <span>
          <strong>
            {requirement.title}
            {requirement.required ? '\u00a0*' : ''}
          </strong>{' '}
          <small>{filename || 'Нажмите, чтобы прикрепить документ'}</small>
        </span>
        <span className="document-upload-status" aria-hidden="true">
          {filename ? '✓' : '+'}
        </span>
      </button>
      <p id={`${id}-description`} className="document-upload-description">
        {requirement.description || 'PDF, PNG или JPEG; до 10 МБ.'}
      </p>
      {filename && (
        <button
          className="text-button document-remove"
          type="button"
          aria-label={`Удалить документ ${filename}`}
          onClick={() => {
            setError('')
            selectedFiles.delete(requirement.id)
            onChange()
          }}
        >
          Удалить
        </button>
      )}
      {(error || missing) && (
        <p className="error-text" role="alert">
          {error || `Прикрепите документ «${requirement.title}».`}
        </p>
      )}
    </div>
  )
}

export function legacyDocumentFiles(documents: string[] = []) {
  const files: Record<string, string> = {}
  for (const filename of documents) {
    if (/паспорт/i.test(filename)) files.passport = filename
    else if (/отч[её]т|финанс|бухгалтер/i.test(filename)) files['annual-report'] = filename
    else if (/анкет.*физ/i.test(filename)) files['person-form'] = filename
    else if (/анкет/i.test(filename)) files['company-form'] = filename
  }
  return files
}
