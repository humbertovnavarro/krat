package jobs

const serializeJobStatement = `
	INSERT INTO TABLE jobs (
		uuid, 
		status,
		executable,
		args,
		expires,
		stopsAt,
		startsAt,
		output
	)
	VALUES (
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?
	)
`
const createJobTableStatement = `
CREATE TABLE IF NOT EXISTS jobs (
	uuid integer not null primary key,
	status text, executable text,
	args text, expires integer,
	stopsAt integer,
	startsAt integer,
	output blob
)
`
