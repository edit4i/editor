<script lang="ts">
  import { Send, MoreVertical, Trash2, Settings, BotIcon } from 'lucide-svelte';
  import Button from '../components/Button.svelte';
  import Input from '../components/Input.svelte';
  import Select from '../components/Select.svelte';
  import DropdownMenu from '../components/DropdownMenu.svelte';
  import { onMount, onDestroy } from 'svelte';
  import { setKeyboardContext } from '@/stores/keyboardStore';

  export let collapsed = false;

  interface Message {
    role: 'assistant' | 'user';
    content: string;
  }

  let messages: Message[] = [
    { role: 'assistant', content: 'Hello! How can I assist you with your code today?' },
    { role: 'user', content: 'Can you explain how to use React hooks?' },
  ];

  let newMessage = '';
  let selectedModel = 'GPT-4';
  let showMoreOptions = false;
  const modelOptions = ['GPT-4', 'GPT-3.5-turbo', 'Codex', 'GPT-3', 'GPT-2', 'GPT-1'];

  function sendMessage() {
    if (newMessage.trim()) {
      messages = [...messages, { role: 'user', content: newMessage.trim() }];
      newMessage = '';
      // Here you would typically send the message to an AI service and await a response
    }
  }

  function handleKeyPress(event: KeyboardEvent) {
    if (event.key === 'Enter') {
      sendMessage();
    }
  }

  onMount(() => {
    if (!collapsed) {
      setKeyboardContext('aiAssistant');
    }
  });

  onDestroy(() => {
    setKeyboardContext('global');
  });

  $: if (collapsed) {
    setKeyboardContext('global');
  } else {
    setKeyboardContext('aiAssistant');
  }
</script>

{#if !collapsed}
  <div class="bg-gray-900 h-full w-full flex flex-col overflow-hidden border-l border-gray-800">
    <div class="h-[35px] flex items-center justify-between px-4 border-b border-gray-800">
      <div class="flex space-x-2 text-sm font-medium flex gap-2">
        <BotIcon size={16} class="mt-[.1rem]" />
        <span>AI Assistant</span>
      </div>
      <div class="flex items-center space-x-2">
        <div class="w-40">
          <Select
            bind:value={selectedModel}
            options={modelOptions}
            variant="compact"
          />
        </div>
        <Button
          variant="ghost"
          size="sm"
          icon={MoreVertical}
          title="More Options"
          on:click={() => showMoreOptions = !showMoreOptions}
        />
        <DropdownMenu
          show={showMoreOptions}
          onClose={() => showMoreOptions = false}
          position="bottom"
          items={[
            {
              icon: Settings,
              label: 'Settings',
              onClick: () => {
                // TODO: Implement settings dialog
              }
            },
            {
              icon: Trash2,
              label: 'Clear Chat',
              onClick: () => messages = []
            }
          ]}
        />
      </div>
    </div>

    <div class="flex-grow overflow-y-auto p-4">
      {#each messages as message, index}
        <div class="mb-4 {message.role === 'user' ? 'text-right' : ''}">
          <div class="inline-block bg-gray-800 rounded-sm p-2 max-w-[80%] text-gray-300">
            {message.content}
          </div>
        </div>
      {/each}
    </div>

    <div class="p-4 border-t border-gray-800">
      <div class="flex items-start">
        <div class="flex-grow mr-2">
          <Input
            variant="textarea"
            bind:value={newMessage}
            placeholder="Type your message..."
            minRows={1}
            maxRows={5}
            onSubmit={sendMessage}
          />
        </div>
        <Button
          variant="secondary"
          icon={Send}
          on:click={sendMessage}
          title="Send message"
        />
      </div>
    </div>
  </div>
{/if}
