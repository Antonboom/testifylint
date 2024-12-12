package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

const GINKGO_FOCUS_EXIT_CODE = 197
const GINKGO_TIME_FORMAT = "01/02/06 15:04:05.999"

// Report captures information about a Ginkgo test run
type Report struct {
	//SuitePath captures the absolute path to the test suite
	SuitePath string

	//SuiteDescription captures the description string passed to the DSL's RunSpecs() function
	SuiteDescription string

	//SuiteLabels captures any labels attached to the suite by the DSL's RunSpecs() function
	SuiteLabels []string

	//SuiteSucceeded captures the success or failure status of the test run
	//If true, the test run is considered successful.
	//If false, the test run is considered unsuccessful
	SuiteSucceeded bool

	//SuiteHasProgrammaticFocus captures whether the test suite has a test or set of tests that are programmatically focused
	//(i.e an `FIt` or an `FDescribe`
	SuiteHasProgrammaticFocus bool

	//SpecialSuiteFailureReasons may contain special failure reasons
	//For example, a test suite might be considered "failed" even if none of the individual specs
	//have a failure state.  For example, if the user has configured --fail-on-pending the test suite
	//will have failed if there are pending tests even though all non-pending tests may have passed.  In such
	//cases, Ginkgo populates SpecialSuiteFailureReasons with a clear message indicating the reason for the failure.
	//SpecialSuiteFailureReasons is also populated if the test suite is interrupted by the user.
	//Since multiple special failure reasons can occur, this field is a slice.
	SpecialSuiteFailureReasons []string

	//PreRunStats contains a set of stats captured before the test run begins.  This is primarily used
	//by Ginkgo's reporter to tell the user how many specs are in the current suite (PreRunStats.TotalSpecs)
	//and how many it intends to run (PreRunStats.SpecsThatWillRun) after applying any relevant focus or skip filters.
	PreRunStats PreRunStats

	//StartTime and EndTime capture the start and end time of the test run
	StartTime time.Time
	EndTime   time.Time

	//RunTime captures the duration of the test run
	RunTime time.Duration

	//SuiteConfig captures the Ginkgo configuration governing this test run
	//SuiteConfig includes information necessary for reproducing an identical test run,
	//such as the random seed and any filters applied during the test run
	SuiteConfig SuiteConfig

	//SpecReports is a list of all SpecReports generated by this test run
	//It is empty when the SuiteReport is provided to ReportBeforeSuite
	SpecReports SpecReports
}

// PreRunStats contains a set of stats captured before the test run begins.  This is primarily used
// by Ginkgo's reporter to tell the user how many specs are in the current suite (PreRunStats.TotalSpecs)
// and how many it intends to run (PreRunStats.SpecsThatWillRun) after applying any relevant focus or skip filters.
type PreRunStats struct {
	TotalSpecs       int
	SpecsThatWillRun int
}

// Add is used by Ginkgo's parallel aggregation mechanisms to combine test run reports form individual parallel processes
// to form a complete final report.
func (report Report) Add(other Report) Report {
	report.SuiteSucceeded = report.SuiteSucceeded && other.SuiteSucceeded

	if other.StartTime.Before(report.StartTime) {
		report.StartTime = other.StartTime
	}

	if other.EndTime.After(report.EndTime) {
		report.EndTime = other.EndTime
	}

	specialSuiteFailureReasons := []string{}
	reasonsLookup := map[string]bool{}
	for _, reasons := range [][]string{report.SpecialSuiteFailureReasons, other.SpecialSuiteFailureReasons} {
		for _, reason := range reasons {
			if !reasonsLookup[reason] {
				reasonsLookup[reason] = true
				specialSuiteFailureReasons = append(specialSuiteFailureReasons, reason)
			}
		}
	}
	report.SpecialSuiteFailureReasons = specialSuiteFailureReasons
	report.RunTime = report.EndTime.Sub(report.StartTime)

	reports := make(SpecReports, len(report.SpecReports)+len(other.SpecReports))
	copy(reports, report.SpecReports)
	offset := len(report.SpecReports)
	for i := range other.SpecReports {
		reports[i+offset] = other.SpecReports[i]
	}

	report.SpecReports = reports
	return report
}

// SpecReport captures information about a Ginkgo spec.
type SpecReport struct {
	// ContainerHierarchyTexts is a slice containing the text strings of
	// all Describe/Context/When containers in this spec's hierarchy.
	ContainerHierarchyTexts []string

	// ContainerHierarchyLocations is a slice containing the CodeLocations of
	// all Describe/Context/When containers in this spec's hierarchy.
	ContainerHierarchyLocations []CodeLocation

	// ContainerHierarchyLabels is a slice containing the labels of
	// all Describe/Context/When containers in this spec's hierarchy
	ContainerHierarchyLabels [][]string

	// LeafNodeType, LeadNodeLocation, LeafNodeLabels and LeafNodeText capture the NodeType, CodeLocation, and text
	// of the Ginkgo node being tested (typically an NodeTypeIt node, though this can also be
	// one of the NodeTypesForSuiteLevelNodes node types)
	LeafNodeType     NodeType
	LeafNodeLocation CodeLocation
	LeafNodeLabels   []string
	LeafNodeText     string

	// State captures whether the spec has passed, failed, etc.
	State SpecState

	// IsSerial captures whether the spec has the Serial decorator
	IsSerial bool

	// IsInOrderedContainer captures whether the spec appears in an Ordered container
	IsInOrderedContainer bool

	// StartTime and EndTime capture the start and end time of the spec
	StartTime time.Time
	EndTime   time.Time

	// RunTime captures the duration of the spec
	RunTime time.Duration

	// ParallelProcess captures the parallel process that this spec ran on
	ParallelProcess int

	// RunningInParallel captures whether this spec is part of a suite that ran in parallel
	RunningInParallel bool

	//Failure is populated if a spec has failed, panicked, been interrupted, or skipped by the user (e.g. calling Skip())
	//It includes detailed information about the Failure
	Failure Failure

	// NumAttempts captures the number of times this Spec was run.
	// Flakey specs can be retried with ginkgo --flake-attempts=N or the use of the FlakeAttempts decorator.
	// Repeated specs can be retried with the use of the MustPassRepeatedly decorator
	NumAttempts int

	// MaxFlakeAttempts captures whether the spec has been retried with ginkgo --flake-attempts=N or the use of the FlakeAttempts decorator.
	MaxFlakeAttempts int

	// MaxMustPassRepeatedly captures whether the spec has the MustPassRepeatedly decorator
	MaxMustPassRepeatedly int

	// CapturedGinkgoWriterOutput contains text printed to the GinkgoWriter
	CapturedGinkgoWriterOutput string

	// CapturedStdOutErr contains text printed to stdout/stderr (when running in parallel)
	// This is always empty when running in series or calling CurrentSpecReport()
	// It is used internally by Ginkgo's reporter
	CapturedStdOutErr string

	// ReportEntries contains any reports added via `AddReportEntry`
	ReportEntries ReportEntries

	// ProgressReports contains any progress reports generated during this spec.  These can either be manually triggered, or automatically generated by Ginkgo via the PollProgressAfter() decorator
	ProgressReports []ProgressReport

	// AdditionalFailures contains any failures that occurred after the initial spec failure.  These typically occur in cleanup nodes after the initial failure and are only emitted when running in verbose mode.
	AdditionalFailures []AdditionalFailure

	// SpecEvents capture additional events that occur during the spec run
	SpecEvents SpecEvents
}

func (report SpecReport) MarshalJSON() ([]byte, error) {
	//All this to avoid emitting an empty Failure struct in the JSON
	out := struct {
		ContainerHierarchyTexts     []string
		ContainerHierarchyLocations []CodeLocation
		ContainerHierarchyLabels    [][]string
		LeafNodeType                NodeType
		LeafNodeLocation            CodeLocation
		LeafNodeLabels              []string
		LeafNodeText                string
		State                       SpecState
		StartTime                   time.Time
		EndTime                     time.Time
		RunTime                     time.Duration
		ParallelProcess             int
		Failure                     *Failure `json:",omitempty"`
		NumAttempts                 int
		MaxFlakeAttempts            int
		MaxMustPassRepeatedly       int
		CapturedGinkgoWriterOutput  string              `json:",omitempty"`
		CapturedStdOutErr           string              `json:",omitempty"`
		ReportEntries               ReportEntries       `json:",omitempty"`
		ProgressReports             []ProgressReport    `json:",omitempty"`
		AdditionalFailures          []AdditionalFailure `json:",omitempty"`
		SpecEvents                  SpecEvents          `json:",omitempty"`
	}{
		ContainerHierarchyTexts:     report.ContainerHierarchyTexts,
		ContainerHierarchyLocations: report.ContainerHierarchyLocations,
		ContainerHierarchyLabels:    report.ContainerHierarchyLabels,
		LeafNodeType:                report.LeafNodeType,
		LeafNodeLocation:            report.LeafNodeLocation,
		LeafNodeLabels:              report.LeafNodeLabels,
		LeafNodeText:                report.LeafNodeText,
		State:                       report.State,
		StartTime:                   report.StartTime,
		EndTime:                     report.EndTime,
		RunTime:                     report.RunTime,
		ParallelProcess:             report.ParallelProcess,
		Failure:                     nil,
		ReportEntries:               nil,
		NumAttempts:                 report.NumAttempts,
		MaxFlakeAttempts:            report.MaxFlakeAttempts,
		MaxMustPassRepeatedly:       report.MaxMustPassRepeatedly,
		CapturedGinkgoWriterOutput:  report.CapturedGinkgoWriterOutput,
		CapturedStdOutErr:           report.CapturedStdOutErr,
	}

	if !report.Failure.IsZero() {
		out.Failure = &(report.Failure)
	}
	if len(report.ReportEntries) > 0 {
		out.ReportEntries = report.ReportEntries
	}
	if len(report.ProgressReports) > 0 {
		out.ProgressReports = report.ProgressReports
	}
	if len(report.AdditionalFailures) > 0 {
		out.AdditionalFailures = report.AdditionalFailures
	}
	if len(report.SpecEvents) > 0 {
		out.SpecEvents = report.SpecEvents
	}

	return json.Marshal(out)
}

// CombinedOutput returns a single string representation of both CapturedStdOutErr and CapturedGinkgoWriterOutput
// Note that both are empty when using CurrentSpecReport() so CurrentSpecReport().CombinedOutput() will always be empty.
// CombinedOutput() is used internally by Ginkgo's reporter.
func (report SpecReport) CombinedOutput() string {
	if report.CapturedStdOutErr == "" {
		return report.CapturedGinkgoWriterOutput
	}
	if report.CapturedGinkgoWriterOutput == "" {
		return report.CapturedStdOutErr
	}
	return report.CapturedStdOutErr + "\n" + report.CapturedGinkgoWriterOutput
}

// Failed returns true if report.State is one of the SpecStateFailureStates
// (SpecStateFailed, SpecStatePanicked, SpecStateinterrupted, SpecStateAborted)
func (report SpecReport) Failed() bool {
	return report.State.Is(SpecStateFailureStates)
}

// FullText returns a concatenation of all the report.ContainerHierarchyTexts and report.LeafNodeText
func (report SpecReport) FullText() string {
	texts := []string{}
	texts = append(texts, report.ContainerHierarchyTexts...)
	if report.LeafNodeText != "" {
		texts = append(texts, report.LeafNodeText)
	}
	return strings.Join(texts, " ")
}

// Labels returns a deduped set of all the spec's Labels.
func (report SpecReport) Labels() []string {
	out := []string{}
	seen := map[string]bool{}
	for _, labels := range report.ContainerHierarchyLabels {
		for _, label := range labels {
			if !seen[label] {
				seen[label] = true
				out = append(out, label)
			}
		}
	}
	for _, label := range report.LeafNodeLabels {
		if !seen[label] {
			seen[label] = true
			out = append(out, label)
		}
	}

	return out
}

// MatchesLabelFilter returns true if the spec satisfies the passed in label filter query
func (report SpecReport) MatchesLabelFilter(query string) (bool, error) {
	filter, err := ParseLabelFilter(query)
	if err != nil {
		return false, err
	}
	return filter(report.Labels()), nil
}

// FileName() returns the name of the file containing the spec
func (report SpecReport) FileName() string {
	return report.LeafNodeLocation.FileName
}

// LineNumber() returns the line number of the leaf node
func (report SpecReport) LineNumber() int {
	return report.LeafNodeLocation.LineNumber
}

// FailureMessage() returns the failure message (or empty string if the test hasn't failed)
func (report SpecReport) FailureMessage() string {
	return report.Failure.Message
}

// FailureLocation() returns the location of the failure (or an empty CodeLocation if the test hasn't failed)
func (report SpecReport) FailureLocation() CodeLocation {
	return report.Failure.Location
}

// Timeline() returns a timeline view of the report
func (report SpecReport) Timeline() Timeline {
	timeline := Timeline{}
	if !report.Failure.IsZero() {
		timeline = append(timeline, report.Failure)
		if report.Failure.AdditionalFailure != nil {
			timeline = append(timeline, *(report.Failure.AdditionalFailure))
		}
	}
	for _, additionalFailure := range report.AdditionalFailures {
		timeline = append(timeline, additionalFailure)
	}
	for _, reportEntry := range report.ReportEntries {
		timeline = append(timeline, reportEntry)
	}
	for _, progressReport := range report.ProgressReports {
		timeline = append(timeline, progressReport)
	}
	for _, specEvent := range report.SpecEvents {
		timeline = append(timeline, specEvent)
	}
	sort.Sort(timeline)
	return timeline
}

type SpecReports []SpecReport

// WithLeafNodeType returns the subset of SpecReports with LeafNodeType matching one of the requested NodeTypes
func (reports SpecReports) WithLeafNodeType(nodeTypes NodeType) SpecReports {
	count := 0
	for i := range reports {
		if reports[i].LeafNodeType.Is(nodeTypes) {
			count++
		}
	}

	out := make(SpecReports, count)
	j := 0
	for i := range reports {
		if reports[i].LeafNodeType.Is(nodeTypes) {
			out[j] = reports[i]
			j++
		}
	}
	return out
}

// WithState returns the subset of SpecReports with State matching one of the requested SpecStates
func (reports SpecReports) WithState(states SpecState) SpecReports {
	count := 0
	for i := range reports {
		if reports[i].State.Is(states) {
			count++
		}
	}

	out, j := make(SpecReports, count), 0
	for i := range reports {
		if reports[i].State.Is(states) {
			out[j] = reports[i]
			j++
		}
	}
	return out
}

// CountWithState returns the number of SpecReports with State matching one of the requested SpecStates
func (reports SpecReports) CountWithState(states SpecState) int {
	n := 0
	for i := range reports {
		if reports[i].State.Is(states) {
			n += 1
		}
	}
	return n
}

// If the Spec passes, CountOfFlakedSpecs returns the number of SpecReports that failed after multiple attempts.
func (reports SpecReports) CountOfFlakedSpecs() int {
	n := 0
	for i := range reports {
		if reports[i].MaxFlakeAttempts > 1 && reports[i].State.Is(SpecStatePassed) && reports[i].NumAttempts > 1 {
			n += 1
		}
	}
	return n
}

// If the Spec fails, CountOfRepeatedSpecs returns the number of SpecReports that passed after multiple attempts
func (reports SpecReports) CountOfRepeatedSpecs() int {
	n := 0
	for i := range reports {
		if reports[i].MaxMustPassRepeatedly > 1 && reports[i].State.Is(SpecStateFailureStates) && reports[i].NumAttempts > 1 {
			n += 1
		}
	}
	return n
}

// TimelineLocation captures the location of an event in the spec's timeline
type TimelineLocation struct {
	//Offset is the offset (in bytes) of the event relative to the GinkgoWriter stream
	Offset int `json:",omitempty"`

	//Order is the order of the event with respect to other events.  The absolute value of Order
	//is irrelevant.  All that matters is that an event with a lower Order occurs before ane vent with a higher Order
	Order int `json:",omitempty"`

	Time time.Time
}

// TimelineEvent represent an event on the timeline
// consumers of Timeline will need to check the concrete type of each entry to determine how to handle it
type TimelineEvent interface {
	GetTimelineLocation() TimelineLocation
}

type Timeline []TimelineEvent

func (t Timeline) Len() int { return len(t) }
func (t Timeline) Less(i, j int) bool {
	return t[i].GetTimelineLocation().Order < t[j].GetTimelineLocation().Order
}
func (t Timeline) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t Timeline) WithoutHiddenReportEntries() Timeline {
	out := Timeline{}
	for _, event := range t {
		if reportEntry, isReportEntry := event.(ReportEntry); isReportEntry && reportEntry.Visibility == ReportEntryVisibilityNever {
			continue
		}
		out = append(out, event)
	}
	return out
}

func (t Timeline) WithoutVeryVerboseSpecEvents() Timeline {
	out := Timeline{}
	for _, event := range t {
		if specEvent, isSpecEvent := event.(SpecEvent); isSpecEvent && specEvent.IsOnlyVisibleAtVeryVerbose() {
			continue
		}
		out = append(out, event)
	}
	return out
}

// Failure captures failure information for an individual test
type Failure struct {
	// Message - the failure message passed into Fail(...).  When using a matcher library
	// like Gomega, this will contain the failure message generated by Gomega.
	//
	// Message is also populated if the user has called Skip(...).
	Message string

	// Location - the CodeLocation where the failure occurred
	// This CodeLocation will include a fully-populated StackTrace
	Location CodeLocation

	TimelineLocation TimelineLocation

	// ForwardedPanic - if the failure represents a captured panic (i.e. Summary.State == SpecStatePanicked)
	// then ForwardedPanic will be populated with a string representation of the captured panic.
	ForwardedPanic string `json:",omitempty"`

	// FailureNodeContext - one of three contexts describing the node in which the failure occurred:
	// FailureNodeIsLeafNode means the failure occurred in the leaf node of the associated SpecReport. None of the other FailureNode fields will be populated
	// FailureNodeAtTopLevel means the failure occurred in a non-leaf node that is defined at the top-level of the spec (i.e. not in a container). FailureNodeType and FailureNodeLocation will be populated.
	// FailureNodeInContainer means the failure occurred in a non-leaf node that is defined within a container.  FailureNodeType, FailureNodeLocation, and FailureNodeContainerIndex will be populated.
	//
	// FailureNodeType will contain the NodeType of the node in which the failure occurred.
	// FailureNodeLocation will contain the CodeLocation of the node in which the failure occurred.
	// If populated, FailureNodeContainerIndex will be the index into SpecReport.ContainerHierarchyTexts and SpecReport.ContainerHierarchyLocations that represents the parent container of the node in which the failure occurred.
	FailureNodeContext FailureNodeContext `json:",omitempty"`

	FailureNodeType NodeType `json:",omitempty"`

	FailureNodeLocation CodeLocation `json:",omitempty"`

	FailureNodeContainerIndex int `json:",omitempty"`

	//ProgressReport is populated if the spec was interrupted or timed out
	ProgressReport ProgressReport `json:",omitempty"`

	//AdditionalFailure is non-nil if a follow-on failure occurred within the same node after the primary failure.  This only happens when a node has timed out or been interrupted.  In such cases the AdditionalFailure can include information about where/why the spec was stuck.
	AdditionalFailure *AdditionalFailure `json:",omitempty"`
}

func (f Failure) IsZero() bool {
	return f.Message == "" && (f.Location == CodeLocation{})
}

func (f Failure) GetTimelineLocation() TimelineLocation {
	return f.TimelineLocation
}

// FailureNodeContext captures the location context for the node containing the failing line of code
type FailureNodeContext uint

const (
	FailureNodeContextInvalid FailureNodeContext = iota

	FailureNodeIsLeafNode
	FailureNodeAtTopLevel
	FailureNodeInContainer
)

var fncEnumSupport = NewEnumSupport(map[uint]string{
	uint(FailureNodeContextInvalid): "INVALID FAILURE NODE CONTEXT",
	uint(FailureNodeIsLeafNode):     "leaf-node",
	uint(FailureNodeAtTopLevel):     "top-level",
	uint(FailureNodeInContainer):    "in-container",
})

func (fnc FailureNodeContext) String() string {
	return fncEnumSupport.String(uint(fnc))
}
func (fnc *FailureNodeContext) UnmarshalJSON(b []byte) error {
	out, err := fncEnumSupport.UnmarshJSON(b)
	*fnc = FailureNodeContext(out)
	return err
}
func (fnc FailureNodeContext) MarshalJSON() ([]byte, error) {
	return fncEnumSupport.MarshJSON(uint(fnc))
}

// AdditionalFailure capturs any additional failures that occur after the initial failure of a psec
// these typically occur in clean up nodes after the spec has failed.
// We can't simply use Failure as we want to track the SpecState to know what kind of failure this is
type AdditionalFailure struct {
	State   SpecState
	Failure Failure
}

func (f AdditionalFailure) GetTimelineLocation() TimelineLocation {
	return f.Failure.TimelineLocation
}

// SpecState captures the state of a spec
// To determine if a given `state` represents a failure state, use `state.Is(SpecStateFailureStates)`
type SpecState uint

const (
	SpecStateInvalid SpecState = 0

	SpecStatePending SpecState = 1 << iota
	SpecStateSkipped
	SpecStatePassed
	SpecStateFailed
	SpecStateAborted
	SpecStatePanicked
	SpecStateInterrupted
	SpecStateTimedout
)

var ssEnumSupport = NewEnumSupport(map[uint]string{
	uint(SpecStateInvalid):     "INVALID SPEC STATE",
	uint(SpecStatePending):     "pending",
	uint(SpecStateSkipped):     "skipped",
	uint(SpecStatePassed):      "passed",
	uint(SpecStateFailed):      "failed",
	uint(SpecStateAborted):     "aborted",
	uint(SpecStatePanicked):    "panicked",
	uint(SpecStateInterrupted): "interrupted",
	uint(SpecStateTimedout):    "timedout",
})

func (ss SpecState) String() string {
	return ssEnumSupport.String(uint(ss))
}
func (ss SpecState) GomegaString() string {
	return ssEnumSupport.String(uint(ss))
}
func (ss *SpecState) UnmarshalJSON(b []byte) error {
	out, err := ssEnumSupport.UnmarshJSON(b)
	*ss = SpecState(out)
	return err
}
func (ss SpecState) MarshalJSON() ([]byte, error) {
	return ssEnumSupport.MarshJSON(uint(ss))
}

var SpecStateFailureStates = SpecStateFailed | SpecStateTimedout | SpecStateAborted | SpecStatePanicked | SpecStateInterrupted

func (ss SpecState) Is(states SpecState) bool {
	return ss&states != 0
}

// ProgressReport captures the progress of the current spec.  It is, effectively, a structured Ginkgo-aware stack trace
type ProgressReport struct {
	Message           string `json:",omitempty"`
	ParallelProcess   int    `json:",omitempty"`
	RunningInParallel bool   `json:",omitempty"`

	ContainerHierarchyTexts []string     `json:",omitempty"`
	LeafNodeText            string       `json:",omitempty"`
	LeafNodeLocation        CodeLocation `json:",omitempty"`
	SpecStartTime           time.Time    `json:",omitempty"`

	CurrentNodeType      NodeType     `json:",omitempty"`
	CurrentNodeText      string       `json:",omitempty"`
	CurrentNodeLocation  CodeLocation `json:",omitempty"`
	CurrentNodeStartTime time.Time    `json:",omitempty"`

	CurrentStepText      string       `json:",omitempty"`
	CurrentStepLocation  CodeLocation `json:",omitempty"`
	CurrentStepStartTime time.Time    `json:",omitempty"`

	AdditionalReports []string `json:",omitempty"`

	CapturedGinkgoWriterOutput string           `json:",omitempty"`
	TimelineLocation           TimelineLocation `json:",omitempty"`

	Goroutines []Goroutine `json:",omitempty"`
}

func (pr ProgressReport) IsZero() bool {
	return pr.CurrentNodeType == NodeTypeInvalid
}

func (pr ProgressReport) Time() time.Time {
	return pr.TimelineLocation.Time
}

func (pr ProgressReport) SpecGoroutine() Goroutine {
	for _, goroutine := range pr.Goroutines {
		if goroutine.IsSpecGoroutine {
			return goroutine
		}
	}
	return Goroutine{}
}

func (pr ProgressReport) HighlightedGoroutines() []Goroutine {
	out := []Goroutine{}
	for _, goroutine := range pr.Goroutines {
		if goroutine.IsSpecGoroutine || !goroutine.HasHighlights() {
			continue
		}
		out = append(out, goroutine)
	}
	return out
}

func (pr ProgressReport) OtherGoroutines() []Goroutine {
	out := []Goroutine{}
	for _, goroutine := range pr.Goroutines {
		if goroutine.IsSpecGoroutine || goroutine.HasHighlights() {
			continue
		}
		out = append(out, goroutine)
	}
	return out
}

func (pr ProgressReport) WithoutCapturedGinkgoWriterOutput() ProgressReport {
	out := pr
	out.CapturedGinkgoWriterOutput = ""
	return out
}

func (pr ProgressReport) WithoutOtherGoroutines() ProgressReport {
	out := pr
	filteredGoroutines := []Goroutine{}
	for _, goroutine := range pr.Goroutines {
		if goroutine.IsSpecGoroutine || goroutine.HasHighlights() {
			filteredGoroutines = append(filteredGoroutines, goroutine)
		}
	}
	out.Goroutines = filteredGoroutines
	return out
}

func (pr ProgressReport) GetTimelineLocation() TimelineLocation {
	return pr.TimelineLocation
}

type Goroutine struct {
	ID              uint64
	State           string
	Stack           []FunctionCall
	IsSpecGoroutine bool
}

func (g Goroutine) IsZero() bool {
	return g.ID == 0
}

func (g Goroutine) HasHighlights() bool {
	for _, fc := range g.Stack {
		if fc.Highlight {
			return true
		}
	}

	return false
}

type FunctionCall struct {
	Function        string
	Filename        string
	Line            int
	Highlight       bool     `json:",omitempty"`
	Source          []string `json:",omitempty"`
	SourceHighlight int      `json:",omitempty"`
}

// NodeType captures the type of a given Ginkgo Node
type NodeType uint

const (
	NodeTypeInvalid NodeType = 0

	NodeTypeContainer NodeType = 1 << iota
	NodeTypeIt

	NodeTypeBeforeEach
	NodeTypeJustBeforeEach
	NodeTypeAfterEach
	NodeTypeJustAfterEach

	NodeTypeBeforeAll
	NodeTypeAfterAll

	NodeTypeBeforeSuite
	NodeTypeSynchronizedBeforeSuite
	NodeTypeAfterSuite
	NodeTypeSynchronizedAfterSuite

	NodeTypeReportBeforeEach
	NodeTypeReportAfterEach
	NodeTypeReportBeforeSuite
	NodeTypeReportAfterSuite

	NodeTypeCleanupInvalid
	NodeTypeCleanupAfterEach
	NodeTypeCleanupAfterAll
	NodeTypeCleanupAfterSuite
)

var NodeTypesForContainerAndIt = NodeTypeContainer | NodeTypeIt
var NodeTypesForSuiteLevelNodes = NodeTypeBeforeSuite | NodeTypeSynchronizedBeforeSuite | NodeTypeAfterSuite | NodeTypeSynchronizedAfterSuite | NodeTypeReportBeforeSuite | NodeTypeReportAfterSuite | NodeTypeCleanupAfterSuite
var NodeTypesAllowedDuringCleanupInterrupt = NodeTypeAfterEach | NodeTypeJustAfterEach | NodeTypeAfterAll | NodeTypeAfterSuite | NodeTypeSynchronizedAfterSuite | NodeTypeCleanupAfterEach | NodeTypeCleanupAfterAll | NodeTypeCleanupAfterSuite
var NodeTypesAllowedDuringReportInterrupt = NodeTypeReportBeforeEach | NodeTypeReportAfterEach | NodeTypeReportBeforeSuite | NodeTypeReportAfterSuite

var ntEnumSupport = NewEnumSupport(map[uint]string{
	uint(NodeTypeInvalid):                 "INVALID NODE TYPE",
	uint(NodeTypeContainer):               "Container",
	uint(NodeTypeIt):                      "It",
	uint(NodeTypeBeforeEach):              "BeforeEach",
	uint(NodeTypeJustBeforeEach):          "JustBeforeEach",
	uint(NodeTypeAfterEach):               "AfterEach",
	uint(NodeTypeJustAfterEach):           "JustAfterEach",
	uint(NodeTypeBeforeAll):               "BeforeAll",
	uint(NodeTypeAfterAll):                "AfterAll",
	uint(NodeTypeBeforeSuite):             "BeforeSuite",
	uint(NodeTypeSynchronizedBeforeSuite): "SynchronizedBeforeSuite",
	uint(NodeTypeAfterSuite):              "AfterSuite",
	uint(NodeTypeSynchronizedAfterSuite):  "SynchronizedAfterSuite",
	uint(NodeTypeReportBeforeEach):        "ReportBeforeEach",
	uint(NodeTypeReportAfterEach):         "ReportAfterEach",
	uint(NodeTypeReportBeforeSuite):       "ReportBeforeSuite",
	uint(NodeTypeReportAfterSuite):        "ReportAfterSuite",
	uint(NodeTypeCleanupInvalid):          "DeferCleanup",
	uint(NodeTypeCleanupAfterEach):        "DeferCleanup (Each)",
	uint(NodeTypeCleanupAfterAll):         "DeferCleanup (All)",
	uint(NodeTypeCleanupAfterSuite):       "DeferCleanup (Suite)",
})

func (nt NodeType) String() string {
	return ntEnumSupport.String(uint(nt))
}
func (nt *NodeType) UnmarshalJSON(b []byte) error {
	out, err := ntEnumSupport.UnmarshJSON(b)
	*nt = NodeType(out)
	return err
}
func (nt NodeType) MarshalJSON() ([]byte, error) {
	return ntEnumSupport.MarshJSON(uint(nt))
}

func (nt NodeType) Is(nodeTypes NodeType) bool {
	return nt&nodeTypes != 0
}

/*
SpecEvent captures a variety of events that can occur when specs run.  See SpecEventType for the list of available events.
*/
type SpecEvent struct {
	SpecEventType SpecEventType

	CodeLocation     CodeLocation
	TimelineLocation TimelineLocation

	Message  string        `json:",omitempty"`
	Duration time.Duration `json:",omitempty"`
	NodeType NodeType      `json:",omitempty"`
	Attempt  int           `json:",omitempty"`
}

func (se SpecEvent) GetTimelineLocation() TimelineLocation {
	return se.TimelineLocation
}

func (se SpecEvent) IsOnlyVisibleAtVeryVerbose() bool {
	return se.SpecEventType.Is(SpecEventByEnd | SpecEventNodeStart | SpecEventNodeEnd)
}

func (se SpecEvent) GomegaString() string {
	out := &strings.Builder{}
	out.WriteString("[" + se.SpecEventType.String() + " SpecEvent] ")
	if se.Message != "" {
		out.WriteString("Message=")
		out.WriteString(`"` + se.Message + `",`)
	}
	if se.Duration != 0 {
		out.WriteString("Duration=" + se.Duration.String() + ",")
	}
	if se.NodeType != NodeTypeInvalid {
		out.WriteString("NodeType=" + se.NodeType.String() + ",")
	}
	if se.Attempt != 0 {
		out.WriteString(fmt.Sprintf("Attempt=%d", se.Attempt) + ",")
	}
	out.WriteString("CL=" + se.CodeLocation.String() + ",")
	out.WriteString(fmt.Sprintf("TL.Offset=%d", se.TimelineLocation.Offset))

	return out.String()
}

type SpecEvents []SpecEvent

func (se SpecEvents) WithType(seType SpecEventType) SpecEvents {
	out := SpecEvents{}
	for _, event := range se {
		if event.SpecEventType.Is(seType) {
			out = append(out, event)
		}
	}
	return out
}

type SpecEventType uint

const (
	SpecEventInvalid SpecEventType = 0

	SpecEventByStart SpecEventType = 1 << iota
	SpecEventByEnd
	SpecEventNodeStart
	SpecEventNodeEnd
	SpecEventSpecRepeat
	SpecEventSpecRetry
)

var seEnumSupport = NewEnumSupport(map[uint]string{
	uint(SpecEventInvalid):    "INVALID SPEC EVENT",
	uint(SpecEventByStart):    "By",
	uint(SpecEventByEnd):      "By (End)",
	uint(SpecEventNodeStart):  "Node",
	uint(SpecEventNodeEnd):    "Node (End)",
	uint(SpecEventSpecRepeat): "Repeat",
	uint(SpecEventSpecRetry):  "Retry",
})

func (se SpecEventType) String() string {
	return seEnumSupport.String(uint(se))
}
func (se *SpecEventType) UnmarshalJSON(b []byte) error {
	out, err := seEnumSupport.UnmarshJSON(b)
	*se = SpecEventType(out)
	return err
}
func (se SpecEventType) MarshalJSON() ([]byte, error) {
	return seEnumSupport.MarshJSON(uint(se))
}

func (se SpecEventType) Is(specEventTypes SpecEventType) bool {
	return se&specEventTypes != 0
}
