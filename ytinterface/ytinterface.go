package ytinterface

import (
	"errors"

	"github.com/kkdai/youtube/v2"
)

func GetBestAudioOnlyFormat(fl youtube.FormatList) (*youtube.Format, error) {
	if f := fl.FindByItag(141); f != nil {
		return f, nil
	}
	if f := fl.FindByItag(140); f != nil {
		return f, nil
	}
	if f := fl.FindByItag(139); f != nil {
		return f, nil
	}
	return nil, errors.New("no audio only format found")
}
