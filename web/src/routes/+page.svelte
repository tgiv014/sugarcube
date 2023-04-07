<script lang="ts">
	import { goto } from '$app/navigation';

	function onSubmit(e: SubmitEvent) {
		const formData = new FormData(e.target as HTMLFormElement);
		const formObj = Object.fromEntries(formData.entries());
		fetch('/api/login', {
			method: 'POST',
			body: JSON.stringify(formObj)
		}).then((response) => {
			if (response.status != 200) {
				return;
			}

			goto('/app');
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
	</div>
</div>
