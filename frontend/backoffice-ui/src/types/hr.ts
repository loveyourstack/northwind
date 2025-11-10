
export interface EmployeeInput {
  address: string
  city: string
  country_fk: number | undefined
  date_of_birth: Date | undefined
  family_name: string
  given_name: string
  hire_date: Date | undefined
  home_phone?: string
  job_title: string
  name: string
  notes: string
  postal_code: string
  reports_to_fk: number | undefined
  state?: string
  title: string
}
export interface Employee extends EmployeeInput {
  id: number
  age: number
  country: string
  country_iso2: string
  reports_to: string
  territory_count: number
}
export function NewEmployee(): Employee {
  return  {
    address: '',
    city: '',
    country_fk: undefined,
    date_of_birth: undefined,
    family_name: '',
    given_name: '',
    hire_date: undefined,
    home_phone: '',
    job_title: '',
    name: '',
    notes: '',
    postal_code: '',
    reports_to_fk: undefined,
    state: '',
    title: '',
    id: 0,
    age: 0,
    country: '',
    country_iso2: '',
    reports_to: '',
    territory_count: 0
  }
}
export function GetEmployeeInputFromItem(item: Employee): EmployeeInput {
  return  {
    address: item.address,
    city: item.city,
    country_fk: item.country_fk,
    date_of_birth: item.date_of_birth,
    family_name: item.family_name,
    given_name: item.given_name,
    hire_date: item.hire_date,
    home_phone: item.home_phone,
    job_title: item.job_title,
    name: item.name,
    notes: item.notes,
    postal_code: item.postal_code,
    reports_to_fk: item.reports_to_fk,
    state: item.state,
    title: item.title,
  }
}

// ------------------------------------------------------------------------------------------------------

export interface MeetingScheduleInput {
  day: string
  frequency: string
  name: string
  scheduled_time: string | undefined
}
export interface MeetingSchedule extends MeetingScheduleInput {
  id: number
}
export function NewMeetingSchedule(): MeetingSchedule {
  return  {
    day: 'None',
    frequency: '',
    name: '',
    scheduled_time: undefined,
    id: 0,
  }
}
export function GetMeetingScheduleInputFromItem(item: MeetingSchedule): MeetingScheduleInput {
  return  {
    day: item.day,
    frequency: item.frequency,
    name: item.name,
    scheduled_time: item.scheduled_time,
  }
}