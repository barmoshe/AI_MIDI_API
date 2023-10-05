from temporalio.client import Client
from temporalio.worker import Worker
from activities import ValidateMIDIText, GenerateMIDIFile
import asyncio


async def main():
    client = await Client.connect("localhost:7233")
    worker = Worker(
        client,
        task_queue="pyhton-worker",
        activities=[ValidateMIDIText, GenerateMIDIFile],
    )
    await worker.run()


if __name__ == "__main__":
    asyncio.run(main())