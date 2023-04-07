<script lang="ts">
	import { goto } from '$app/navigation';
	import { login, status } from '../lib/stores';

	let error = '';

	$: if ($status) {
		if ($status.sessionValid) {
			goto('/app');
		} else if (!$status.setupCompleted) {
			goto('/setup');
		}
	}

	function onSubmit(e: SubmitEvent) {
		const formData = new FormData(e.target as HTMLFormElement);
		const formObj = Object.fromEntries(formData.entries());
		login({
			password: formObj['password'].toString()
		})
			.then(() => {
				goto('/app');
			})
			.catch((reason) => {
				error = reason;
			});
	}
</script>

<div class="flex min-h-screen items-center justify-center">
	<div class="border border-stone-900 p-8">
		<h1 class="mb-6 font-mono text-4xl font-thin italic">Welcome back.</h1>
		<p />
		<form on:submit|preventDefault={onSubmit}>
			<label class="flex flex-col gap-4 font-mono">
				Please enter your password<input
					name="password"
					id="password"
					type="password"
					autocomplete="off"
					placeholder="right"
					class="rounded-full border border-stone-900 px-4 text-xl"
				/></label
			>
		</form>
		{#if error}
			<p>{error}</p>
		{/if}
	</div>
</div>
