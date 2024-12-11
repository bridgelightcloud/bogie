# Bogie Transit Tracker
Known simply as "Bogie", this project is a system for tracking and analyzing transit data.

### Models

#### Event
| Field | Type | DynamoDB Name | Description |
|-------|------|---------------|-------------|
| `Id` | `uuid.UUID` |  `id` | UUID |
| `Type` | `string` | `t` | Doctype, "event", "user", etc for identificaiton in the db table |
| `Status` | `string` | `s` | Active/Inactive |
| `CreatedAt` | `*time.Time` | `ca` | Time record was created |
| `UpdatedAt` | `*time.Time` | `ua` | Time record was updated |
| `User` | `uuid.UUID` | `uid` | ID of the user who created the event |
| `Agency` | `string` | `a` | Matches [`agency.agency_name`](#gtfs-parser) |
| `Route` | `string` | `r` | Matches [`routes.route_id`](#gtfs-parser) |
| `Trip` | `string` | `tr` | Matches [`trips.trip_id`](#gtfs-parser) |
| `UnitID` | `string` | `u` | ID of the vehicle ridden (to be changed to `VehicleID`, `v`) |
| `UnitCount` | `*int` | `uc` | Number of units (vehicles) in a compliment |
| `UnitPosition` | `*int` | `up` | Position of the ridden vehicle |
| `DepartureStop` | `string` | `ds` | Stop where event was started. Matches [`stops.stop_id`](#gtfs-parser) |
| `ArrivalStop` | `string` | `as` | Stop where event was ended. Matches [`stops.stop_id`](#gtfs-parser) |
| `DepartureTime` | `*time.Time` | `dt` | Time vehicle left departure stop |
| `ArrivalTime` | `*time.Time` | `at` | Time vehicle arrived at arrival stop |
| `Notes` | `[]string` | `n` | String notes associated with this event |

### DB

### API

## Sub Projects

### CSVMUM
Marshal and unmarshal CSV files to and from Go structs, using reflection, tags, and custom parsers

[README](./pkg/csvmum/README.md)

### GTFS Parser
Read GTFS zip files and parse the data

[README](./pkg/gtfs/README.md)