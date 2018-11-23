package blockvolmanager

import "strconv"

// OptFunc configures optional parameters to be used for a block volume operation like create,delete
type OptFunc func(*Options)

// Options represents optional parameters for a block volume operation like create,delete
type Options struct {
	// optional params for delete req
	forceDelete   bool
	unlinkStorage bool

	// optional params for a create req
	auth               bool
	fullPrealloc       bool
	storage            string
	ha                 int
	ringBufferSizeInMB uint64
}

// applyOpts applies configured optional parameters on Options
func (op *Options) applyOpts(optFuncs ...OptFunc) {
	for _, optFunc := range optFuncs {
		optFunc(op)
	}
}

// prepareArgs returns optional args to be used for a block volume operation
func (op *Options) prepareArgs() (args []string) {
	if op.ha > 0 {
		args = append(args, "ha", strconv.Itoa(op.ha))
	}

	if op.fullPrealloc {
		args = append(args, "prealloc", "full")
	}

	if op.auth {
		args = append(args, "auth", "enable")
	}

	if op.ringBufferSizeInMB > 0 {
		args = append(args, "ring-buffer", strconv.FormatUint(op.ringBufferSizeInMB, 10))
	}

	if op.storage != "" {
		args = append(args, "storage", op.storage)
	}

	if op.unlinkStorage {
		args = append(args, "unlink-storage", "yes")
	}

	if op.forceDelete {
		args = append(args, "force")
	}
	return
}

// WithHaCount configures haCount for block creation
func WithHaCount(count int) OptFunc {
	return func(options *Options) {
		options.ha = count
	}
}

// WithAuthEnabled enables auth for block creation
func WithAuthEnabled(options *Options) {
	options.auth = true
}

// WithFullPrealloc sets "prealloc" args as "full"
func WithFullPrealloc(options *Options) {
	options.fullPrealloc = true
}

// WithStorage configures storage args for block-vol creation
func WithStorage(storage string) OptFunc {
	return func(options *Options) {
		options.storage = storage
	}
}

// WithRingBufferSizeInMB configures ring-buffer args (size should in MB units)
func WithRingBufferSizeInMB(size uint64) OptFunc {
	return func(options *Options) {
		options.ringBufferSizeInMB = size
	}
}

// WithForceDelete configures force args in a block delete req
func WithForceDelete(options *Options) {
	options.forceDelete = true
}

// WithUnlinkStorage configures unlink-storage args in block delete req
func WithUnlinkStorage(options *Options) {
	options.forceDelete = true
}
