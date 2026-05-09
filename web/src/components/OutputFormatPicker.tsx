import { getOutputFormatLabel, OUTPUT_FORMATS } from '../lib/ratios'

export function OutputFormatPicker({ value, onChange }: { value: string; onChange: (format: string) => void }) {
  return (
    <div className="format-list" role="radiogroup" aria-label="输出格式">
      {OUTPUT_FORMATS.map((format) => (
        <button
          key={format}
          type="button"
          className={`format-btn ${format === value ? 'active' : ''}`}
          onClick={() => onChange(format)}
          aria-checked={format === value}
          role="radio"
          title={`输出 ${getOutputFormatLabel(format)} 格式`}
        >
          {getOutputFormatLabel(format)}
        </button>
      ))}
    </div>
  )
}
