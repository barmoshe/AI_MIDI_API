from temporalio import activity

@activity.defn(name="ValidateMIDIText")
async def ValidateMIDIText(midi_text) -> bool:
    return true

@activity.defn(name="GenerateMIDIFile")
async def GenerateMIDIFile(midi_text) -> str:
    return "path/to/midi/file.mid"
    
    
