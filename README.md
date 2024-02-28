# Notepad Cache Spy

Quick'n'dirty implementation to read cache files of Notepad. Inspired by this [YouTube video](https://youtu.be/zSSBbv2fc2s).

## Installation
```bash
CGO_ENABLED=0 go build .
```

## Usage
```bash
npcs.exe # to scan all files of the current user OR
npcs.exe example.bin # to read a specific file
```

## Example Output
### Unsaved Buffer
```bash
$ npcs.exe C:\Users\Infosec\AppData\Local\Packages\Microsoft.WindowsNotepad_8wekyb3d8bbwe\LocalState\TabState\04775658-4985-43fa-8ddd-301d202ca61c.bin
[i] Type:    0x4E50000001
[i] File:    <memory>
[i] Unknown: e10201000000e102
[i] Content:
    com.facebook.services
    com.facebook.system
    com.facebook.appmanager
```

### Saved Buffer
```bash
$ npcs.exe C:\Users\Infosec\AppData\Local\Packages\Microsoft.WindowsNotepad_8wekyb3d8bbwe\LocalState\TabState\99f37c49-4495-4fd4-9e61-9d1128492e84.bin
[i] Type:    0x4E500001
[i] File:    C:\Data\some_text.txt
[i] Unknown: 84030501978b8df6face9aed01921b9297e47bf2b060dc4b2e12fff3c4036cf030f14390f66f14f6e14d5c47ac0001f502f50201000000f502
[i] Content:
    Hello World, this is some random text. Note that the `Unknown` value listed above will probably not match this text. 
    Since it's not clear what these unknown bytes do, you might encounter some garbage text at the content start. 
    Feel free to submit PRs to fix edge cases.

```
