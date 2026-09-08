package searcher

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"serverigh/internal/filesystem/fileinfo"
)

func (s Service) enrichResults(
	ctx context.Context,
	results []Result,
	limit int,
) ([]Result, error) {
	enriched := make([]Result, 0, min(limit, len(results)))
	for _, result := range results {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}

		enrichedResult, ok, err := s.enrichResult(result)
		if err != nil {
			return nil, err
		}

		if ok {
			enriched = append(enriched, enrichedResult)
		}

		if len(enriched) >= limit {
			return enriched, nil
		}
	}

	return enriched, nil
}

func (s Service) enrichResult(result Result) (Result, bool, error) {
	filename := filepath.Join(
		s.root,
		filepath.FromSlash(strings.TrimPrefix(result.Path, "/")),
	)

	info, err := os.Lstat(filename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) || errors.Is(err, fs.ErrPermission) {
			return Result{}, false, nil
		}

		return Result{}, false, wrapFileError("inspect search entry", err)
	}

	createdTime, createdKnown := fileinfo.CreationTime(filename)
	result.Kind = fileinfo.Kind(info)
	result.Size = info.Size()
	result.SizeLabel = fileinfo.FormatSize(info.Size())
	result.RelativePath = strings.TrimPrefix(result.Path, "/")
	result.AbsolutePath = filename
	result.Mode = info.Mode().String()
	result.ModTimeLabel = fileinfo.FormatTime(info.ModTime())
	result.ModTimeValue = fileinfo.FormatTimeValue(info.ModTime())
	result.CreatedLabel = fileinfo.FormatOptionalTime(createdTime, createdKnown)
	result.CreatedValue = fileinfo.FormatOptionalTimeValue(
		createdTime,
		createdKnown,
	)
	result.CreatedKnown = createdKnown
	result.IsDir = info.IsDir()

	return result, true, nil
}
