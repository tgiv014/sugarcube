<script lang="ts">
	import { settings, updateSettings, type SettingsUpdate } from '$lib/stores';
	let error = '';
	let dexcomUsername = $settings.dexcomUsername;
	let dexcomPassword = '';

	function onSubmit() {
		updateSettings({
			dexcomUsername,
			dexcomPassword
		})
			.then(() => {
				dexcomPassword = '';
			})
			.catch((reason) => {
				error = reason;
			});
	}

	function resetForm() {
		dexcomUsername = $settings.dexcomUsername;
		dexcomPassword = '';
	}
</script>

<div class="mx-auto max-w-4xl p-8">
	<h1 class="mb-4 text-2xl">Settings</h1>
	<div class="flex flex-col gap-4 rounded-lg bg-stone-700 p-8">
		<h2 class="text-lg">Dexcom</h2>
		<p class="text-sm">
			To get access to glucose readings, I need your Dexcom username and password. This will be used
			to access the Dexcom Share API for realtime readings. This data never leaves your sugarcube
			instance, except when authenticating with Dexcom.
		</p>
		<form on:submit|preventDefault={onSubmit} class="flex flex-col gap-4">
			<div class="ml-8 flex flex-col gap-4">
				<label class="flex items-center gap-8">
					Dexcom Username
					<input
						name="dexcomUsername"
						id="dexcomUsername"
						type="text"
						autocomplete="off"
						placeholder="username"
						bind:value={dexcomUsername}
						class="rounded-full border border-stone-900 bg-stone-100 px-4 py-1 text-xl dark:border-stone-100 dark:bg-stone-900"
					/>
				</label>
				<label class="flex items-center gap-8">
					Dexcom Password
					<input
						name="dexcomPassword"
						id="dexcomPassword"
						type="password"
						autocomplete="off"
						placeholder="password"
						bind:value={dexcomPassword}
						class="rounded-full border border-stone-900 bg-stone-100 px-4 py-1 text-xl dark:border-stone-100 dark:bg-stone-900"
					/>
				</label>
			</div>
			<div class="flex justify-end gap-4">
				<button
					class="rounded-full border-2 border-cyan-600 px-4 py-2"
					on:click|preventDefault={resetForm}
					type="button"
				>
					Cancel
				</button>
				<button type="submit" class="rounded-full bg-cyan-600 px-4 py-2"> Save </button>
			</div>
		</form>
	</div>
</div>
