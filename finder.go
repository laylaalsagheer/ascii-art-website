package main

import (

	"os"
	"strings"
)

func GenerateASCIIArt(text, bannerType string) string {
    // Map banner types to their corresponding file names
    bannerFiles := map[string]string{
        "shadow":     "shadow.txt",
        "standard":   "standard.txt",
        "thinkertoy": "thinkertoy.txt",
    }

    // Validate the selected banner type
    selectedBannerFile, exists := bannerFiles[bannerType]
    if !exists {
        return "Invalid banner type"
    }
    // Open the selected banner text file
    file, err := os.Open(selectedBannerFile)
    if err != nil {
        return "Error opening banner file"
    }
    defer file.Close()

    // Read the content of the selected banner text file
    bannerContent, err := os.ReadFile(selectedBannerFile)
    if err != nil {
        return "Error reading banner file"
    }
    
    bannerLines := strings.Split(string(bannerContent), "\n")

    var result strings.Builder

    // Generate ASCII art by mapping text characters to the banner columns
    letters := []rune(text)

    for i := 0; i < 8; i++ {
        for _, char := range letters {
          //  result.WriteString("  ")
            result.WriteString(bannerLines[ int((char - 32)*9)+1 + i])
        }
        result.WriteString("\n")

    }

    return result.String()
}