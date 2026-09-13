import { useState } from 'react';
import Select from '../Select';
import type { Period, TimeRange } from './model';

const presets: Record<Period, number[]> = {
  Minutes: [1, 5, 10, 15, 30, 45], Hours: [1, 2, 3, 6, 8, 12],
  Days: [1, 2, 3, 4, 5, 6], Weeks: [1, 2, 3], Months: [1, 2, 3, 4, 5, 6],
};
export default function DateTime({ value, onChange }: { value: TimeRange; onChange: (value: TimeRange) => void }) {
  const [open, setOpen] = useState(false);
  const update = (patch: Partial<TimeRange>) => onChange({ ...value, ...patch });
  const label = value.selectedFullTime ? 'FullTime' : value.tab === 'relative'
    ? `${value.selectedRelativeValue} ${value.selectedRelativePeriod}`
    : `${value.startDate} ${value.startTime} - ${value.endDate} ${value.endTime}`;
  return <div className="relative">
    <button type="button" data-cy="date-time-button" aria-expanded={open} onClick={() => setOpen(!open)}>{label}</button>
    {open && <section className="card absolute right-0 z-20 w-max max-w-[90vw]" aria-label="Time range">
      <div className="toolbar" role="tablist" aria-label="Time range type">
        {(['relative', 'absolute'] as const).map(tab => <button type="button" role="tab" aria-selected={value.tab === tab} key={tab} onClick={() => update({ tab })}>{tab}</button>)}
      </div>
      {value.tab === 'relative' ? <div>
        {Object.entries(presets).map(([period, values]) => <div className="toolbar" key={period}><span className="w-20">{period}</span>{values.map(amount => <button type="button" key={amount} aria-label={`${amount} ${period}`} aria-pressed={value.selectedRelativePeriod === period && value.selectedRelativeValue === amount} className={value.selectedRelativePeriod === period && value.selectedRelativeValue === amount ? 'primary' : ''} onClick={() => update({ selectedRelativePeriod: period as Period, selectedRelativeValue: amount })}>{amount}</button>)}</div>)}
        <div className="toolbar"><label>Custom<input aria-label="Relative value" type="number" min="1" value={value.selectedRelativeValue} onChange={event => update({ selectedRelativeValue: Number(event.target.value) })} /></label>
          <label>Period<Select aria-label="Relative period" value={value.selectedRelativePeriod} onValueChange={value => update({ selectedRelativePeriod: value as Period })}>{Object.keys(presets).map(period => <option key={period}>{period}</option>)}</Select></label></div>
      </div> : <div className="grid grid-cols-2 gap-2">
        {(['startDate', 'endDate', 'startTime', 'endTime'] as const).map(field => <label key={field}>{({ startDate: 'Start Date', endDate: 'End Date', startTime: 'Start Time', endTime: 'End Time' })[field]}<input type={field.endsWith('Date') ? 'date' : 'time'} value={value[field]} onChange={event => update({ [field]: event.target.value })} /></label>)}
      </div>}
      <div className="time-range-actions">
        <label><input type="checkbox" checked={value.selectedFullTime} onChange={event => update({ selectedFullTime: event.target.checked })} /> FullTime</label>
        <button type="button" onClick={() => setOpen(false)}>Close</button>
      </div>
    </section>}
  </div>;
}
