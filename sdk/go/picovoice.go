// Copyright 2021 Picovoice Inc.
//
// You may not use this file except in compliance with the license. A copy of the license is
// located in the "LICENSE" file accompanying this source.
//
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing permissions and
// limitations under the License.
//

// Go binding for Picovoice end-to-end platform. Picovoice enables building voice experiences similar to Alexa but
// runs entirely on-device (offline).

// Picovoice detects utterances of a customizable wake word (phrase) within an incoming stream of audio in real-time.
// After detection of wake word, it begins to infer the user's intent from the follow-on spoken command. Upon detection
// of wake word and completion of voice command, it invokes user-provided callbacks to signal these events.

// Picovoice processes incoming audio in consecutive frames. The number of samples per frame is
// `FrameLength`. The incoming audio needs to have a sample rate equal to `SampleRate` and be 16-bit
// linearly-encoded. Picovoice operates on single-channel audio. It uses Porcupine wake word engine for wake word
// detection and Rhino Speech-to-Intent engine for intent inference.

package picovoice

import (
	"fmt"
	"os"

	ppn "github.com/Picovoice/porcupine/binding/go"
	rhn "github.com/Picovoice/rhino/binding/go"
)

// PvStatus descibes error codes returned from native code
type PvStatus int

const (
	SUCCESS          PvStatus = 0
	OUT_OF_MEMORY    PvStatus = 1
	IO_ERROR         PvStatus = 2
	INVALID_ARGUMENT PvStatus = 3
	STOP_ITERATION   PvStatus = 4
	KEY_ERROR        PvStatus = 5
	INVALID_STATE    PvStatus = 6
)

func pvStatusToString(status PvStatus) string {
	switch status {
	case SUCCESS:
		return "SUCCESS"
	case OUT_OF_MEMORY:
		return "OUT_OF_MEMORY"
	case IO_ERROR:
		return "IO_ERROR"
	case INVALID_ARGUMENT:
		return "INVALID_ARGUMENT"
	case STOP_ITERATION:
		return "STOP_ITERATION"
	case KEY_ERROR:
		return "KEY_ERROR"
	case INVALID_STATE:
		return "INVALID_STATE"
	default:
		return fmt.Sprintf("Unknown error code: %d", status)
	}
}

// Callback for when a wake word has been detected
type WakeWordCallbackType func()

// Callback for when Rhino has made an inference
type InferenceCallbackType func(rhn.RhinoInference)

// Picovoice struct
type Picovoice struct {
	// instance of porcupine
	porcupine ppn.Porcupine

	// instance of rhino
	rhino rhn.Rhino

	// only true after init and before delete
	initialized bool

	// true after Porcupine detected wake word
	wakeWordDetected bool

	// Path to Porcupine keyword file (.ppn)
	KeywordPath string

	// Function to be called once the wake word has been detected
	WakeWordCallback WakeWordCallbackType

	// Path to Rhino context file (.rhn)
	ContextPath string

	// Function to be called once Rhino has an inference ready
	InferenceCallback InferenceCallbackType

	// Path to Porcupine model file (.pv)
	PorcupineModelPath string

	// Sensitivity value for detecting keyword. The value should be a number within [0, 1]. A
	// higher sensitivity results in fewer misses at the cost of increasing the false alarm rate.
	PorcupineSensitivity float32

	// Path to Rhino model file (.pv)
	RhinoModelPath string

	// Inference sensitivity. A higher sensitivity value results in
	// fewer misses at the cost of (potentially) increasing the erroneous inference rate.
	// Sensitivity should be a floating-point number within 0 and 1.
	RhinoSensitivity float32

	// Once initialized, stores the source of the Rhino context in YAML format. Shows the list of intents,
	// which expressions map to those intents, as well as slots and their possible values.
	ContextInfo string
}

// Returns a Picovoice stuct with default parameters
func NewPicovoice(keywordPath string,
	wakewordCallback WakeWordCallbackType,
	contextPath string,
	inferenceCallback InferenceCallbackType) Picovoice {
	return Picovoice{
		KeywordPath:       keywordPath,
		WakeWordCallback:  wakewordCallback,
		ContextPath:       contextPath,
		InferenceCallback: inferenceCallback,

		PorcupineSensitivity: 0.5,
		RhinoSensitivity:     0.5,
	}
}

var (
	// Required number of audio samples per frame.
	FrameLength = ppn.FrameLength

	// Required sample rate of input audio
	SampleRate = ppn.SampleRate

	// Version of Porcupine being used
	PorcupineVersion = ppn.Version

	// Version of Rhino being used
	RhinoVersion = rhn.Version

	// Picovoice version
	Version = fmt.Sprintf("1.1.0 (Porcupine v%s) (Rhino v%s)", PorcupineVersion, RhinoVersion)
)

// Init function for Picovoice. Must be called before attempting process.
func (picovoice *Picovoice) Init() error {

	if picovoice.KeywordPath == "" {
		return fmt.Errorf("%s: No valid keyword was provided.", pvStatusToString(INVALID_ARGUMENT))
	}

	if _, err := os.Stat(picovoice.KeywordPath); os.IsNotExist(err) {
		return fmt.Errorf("%s: Keyword file file could not be found at %s", pvStatusToString(INVALID_ARGUMENT), picovoice.KeywordPath)
	}

	if picovoice.ContextPath == "" {
		return fmt.Errorf("%s: No valid context was provided.", pvStatusToString(INVALID_ARGUMENT))
	}

	if _, err := os.Stat(picovoice.ContextPath); os.IsNotExist(err) {
		return fmt.Errorf("%s: Context file could not be found at %s", pvStatusToString(INVALID_ARGUMENT), picovoice.ContextPath)
	}

	if picovoice.InferenceCallback == nil {
		return fmt.Errorf("%s: No InferenceCallback was provided.", pvStatusToString(INVALID_ARGUMENT))
	}

	if ppn.SampleRate != rhn.SampleRate {
		return fmt.Errorf("%s: Pocupine sample rate (%d) was differenct than Rhino sample rate (%d)",
			pvStatusToString(INVALID_ARGUMENT),
			ppn.SampleRate,
			rhn.SampleRate)
	}

	if ppn.FrameLength != rhn.FrameLength {
		return fmt.Errorf("%s: Pocupine frame length (%d) was differenct than Rhino frame length (%d)",
			pvStatusToString(INVALID_ARGUMENT),
			ppn.FrameLength,
			rhn.FrameLength)
	}

	picovoice.porcupine = ppn.Porcupine{
		ModelPath:     picovoice.PorcupineModelPath,
		KeywordPaths:  []string{picovoice.KeywordPath},
		Sensitivities: []float32{0.5},
	}
	err := picovoice.porcupine.Init()
	if err != nil {
		return err
	}

	picovoice.rhino = rhn.Rhino{
		ModelPath:   picovoice.RhinoModelPath,
		ContextPath: picovoice.ContextPath,
		Sensitivity: picovoice.RhinoSensitivity,
	}
	err = picovoice.rhino.Init()
	if err != nil {
		return err
	}
	picovoice.ContextInfo = picovoice.rhino.ContextInfo
	picovoice.initialized = true
	return nil
}

// Releases resouces aquired by Picovoice
func (picovoice *Picovoice) Delete() error {
	porcupineErr := picovoice.porcupine.Delete()
	rhinoErr := picovoice.rhino.Delete()

	if porcupineErr != nil {
		return porcupineErr
	}
	if rhinoErr != nil {
		return rhinoErr
	}

	picovoice.initialized = false
	return nil
}

// Process a frame of pcm audio with the Picovoice platform.
// Invokes user-defined callbacks upon detection of wake word and completion of follow-on command inference
func (picovoice *Picovoice) Process(pcm []int16) error {
	if !picovoice.initialized {
		return fmt.Errorf("Picovoice could not process because it has either not been initialized or has been deleted.")
	}

	if len(pcm) != FrameLength {
		return fmt.Errorf("Input data frame size (%d) does not match required size of %d", len(pcm), FrameLength)
	}

	if !picovoice.wakeWordDetected {
		keywordIndex, err := picovoice.porcupine.Process(pcm)
		if err != nil {
			return err
		}

		if keywordIndex == 0 {
			picovoice.wakeWordDetected = true
			if picovoice.WakeWordCallback != nil {
				picovoice.WakeWordCallback()
			}
		}
	} else {
		isFinalized, err := picovoice.rhino.Process(pcm)
		if err != nil {
			return err
		}
		if isFinalized {
			picovoice.wakeWordDetected = false
			inference, err := picovoice.rhino.GetInference()
			if err != nil {
				return err
			}

			picovoice.InferenceCallback(inference)
		}
	}
	return nil
}
