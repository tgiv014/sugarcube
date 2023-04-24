<script>
	import Footer from './footer.svelte';
	import Header from './header.svelte';

	import { status } from '../../lib/status';
	import { goto } from '$app/navigation';

	$: if ($status) {
		if (!$status.setupCompleted) {
			goto('/setup');
		} else if (!$status.sessionValid) {
			goto('/');
		}
	}
</script>

<div class="flex min-h-screen w-full flex-col">
	<!-- todo: hide header/nav -->
	<Header />
	<div class="flex-1">
		<slot />
	</div>
	<Footer />
</div>
