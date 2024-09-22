package init

import (
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestNewCmdInit(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	cmd := NewCmdInit(&model.GlobalOptions{})

	_, err = cmd.ExecuteC()
	require.NoError(t, err)

	configFilePath := filepath.Join(tempDir, "gencoder.yaml")
	_, err = os.Stat(configFilePath)
	assert.NoError(t, err, "gencoder.yaml should be created")

	templatesDirPath := filepath.Join(tempDir, "templates")
	_, err = os.Stat(templatesDirPath)
	assert.NoError(t, err, "templates directory should be created")

	entityJavaFilePath := filepath.Join(templatesDirPath, "entity.java.hbs")
	_, err = os.Stat(entityJavaFilePath)
	assert.NoError(t, err, "entity.java.hbs should be created")

	javaTypePartialPath := filepath.Join(templatesDirPath, "java_type.partial.hbs")
	_, err = os.Stat(javaTypePartialPath)
	assert.NoError(t, err, "java_type.partial.hbs should be created")
}
