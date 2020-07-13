package config

const (
	MainConfigTemplate = `uuid='{{.UUID}}'
ftdataway='{{.FtGateway}}' # deprecated
log='{{.Log}}'
log_level='{{.LogLevel}}'
config_dir='{{.ConfigDir}}'
http_server_dir='{{.HTTPServerAddr}}'

## Override default hostname, if empty use os.Hostname()
hostname = ""
## If set to true, do no set the "host" tag.
omit_hostname = false

output_file = ""

# ##tell dataway the interval to check datakit alive
#max_post_interval = '1m'

## Global tags can be specified here in key="value" format.
#[global_tags]
# name = 'admin'

[dataway]
	host = "{{.DataWay.Host}}"
	scheme = "{{.DataWay.Scheme}}"
	token = "{{.DataWay.Token}}"
	default_path = "{{.DataWay.DefaultPath}}"

# Configuration for agent
#[agent]
#  ## Default data collection interval for all inputs
#  interval = "10s"
#  ## Rounds collection interval to 'interval'
#  ## ie, if interval="10s" then always collect on :00, :10, :20, etc.
#  round_interval = true

#  ## Telegraf will send metrics to outputs in batches of at most
#  ## metric_batch_size metrics.
#  ## This controls the size of writes that Telegraf sends to output plugins.
#  metric_batch_size = 1000

#  ## Maximum number of unwritten metrics per output.
#  metric_buffer_limit = 100000

#  ## Collection jitter is used to jitter the collection by a random amount.
#  ## Each plugin will sleep for a random time within jitter before collecting.
#  ## This can be used to avoid many plugins querying things like sysfs at the
#  ## same time, which can have a measurable effect on the system.
#  collection_jitter = "0s"

#  ## Default flushing interval for all outputs. Maximum flush_interval will be
#  ## flush_interval + flush_jitter
#  flush_interval = "10s"
#  ## Jitter the flush interval by a random amount. This is primarily to avoid
#  ## large write spikes for users running a large number of telegraf instances.
#  ## ie, a jitter of 5s and interval 10s means flushes will happen every 10-15s
#  flush_jitter = "0s"

#  ## By default or when set to "0s", precision will be set to the same
#  ## timestamp order as the collection interval, with the maximum being 1s.
#  ##   ie, when interval = "10s", precision will be "1s"
#  ##       when interval = "250ms", precision will be "1ms"
#  ## Precision will NOT be used for service inputs. It is up to each individual
#  ## service input to set the timestamp at the appropriate precision.
#  ## Valid time units are "ns", "us" (or "µs"), "ms", "s".
#  precision = "ns"

#  ## Log at debug level.
#  debug = true
#  ## Log only error level messages.
#  # quiet = false

#  ## Log file name, the empty string means to log to stderr.
#  logfile = "/var/log/telegraf/telegraf.log"

#  ## The logfile will be rotated after the time interval specified.  When set
#  ## to 0 no time based rotation is performed.  Logs are rotated only when
#  ## written to, if there is no log activity rotation may be delayed.
#  # logfile_rotation_interval = "0d"

#  ## The logfile will be rotated when it becomes larger than the specified
#  ## size.  When set to 0 no size based rotation is performed.
#  # logfile_rotation_max_size = "0MB"

#  ## Maximum number of rotated archives to keep, any older logs are deleted.
#  ## If set to -1, no archives are removed.
#  # logfile_rotation_max_archives = 5

#  ## Override default hostname, if empty use os.Hostname()
#  hostname = ""
#  ## If set to true, do no set the "host" tag in the telegraf agent.
#  omit_hostname = false

`
)
