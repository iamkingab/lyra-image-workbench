import { useEffect } from 'react'
import { getTask, streamTaskEvents } from '../api/tasks'
import type { Task, TaskEvent } from '../types'

export function useTaskEvents(taskId: string | null, onTask: (task: Task) => void, onEvent?: (event: TaskEvent) => void) {
  useEffect(() => {
    if (!taskId) return
    const controller = new AbortController()
    let stopped = false
    async function run() {
      try {
        onTask(await getTask(taskId!))
        await streamTaskEvents(taskId!, (event) => {
          onEvent?.(event)
          const job = event.data?.job
          if (job) onTask(job)
        }, controller.signal)
      } catch {
        if (!stopped) window.setTimeout(run, 2500)
      }
    }
    void run()
    return () => {
      stopped = true
      controller.abort()
    }
  }, [taskId, onTask, onEvent])
}
