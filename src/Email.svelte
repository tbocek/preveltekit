<script lang="ts">
    let emailPreview = $state({
        subject: "Welcome to LightKit",
        content: "Thank you for trying our framework!",
        template: "default"
    });

    let templates = [
        { id: "default", name: "Default Template" },
        { id: "minimal", name: "Minimal Template" },
        { id: "marketing", name: "Marketing Template" }
    ];

    function updateTemplate(event: Event) {
        const target = event.target as HTMLSelectElement;
        emailPreview.template = target.value;
    }

    // Demo the SSPR capability
    let renderInfo = "Client Rendered";
    if (window?.JSDOM) {
        renderInfo = "Server Pre-Rendered";
    }
</script>

<div class="email-preview">
    <div class="sidebar">
        <h2>Email Builder</h2>
        <p class="render-info">({renderInfo})</p>

        <div class="form-group">
            <label for="subject">Subject</label>
            <input
                    type="text"
                    id="subject"
                    bind:value={emailPreview.subject}
            />
        </div>

        <div class="form-group">
            <label for="template">Template</label>
            <select
                    id="template"
                    bind:value={emailPreview.template}
                    onchange={updateTemplate}
            >
                {#each templates as template}
                    <option value={template.id}>{template.name}</option>
                {/each}
            </select>
        </div>

        <div class="form-group">
            <label for="content">Content</label>
            <textarea
                    id="content"
                    bind:value={emailPreview.content}
                    rows="5"
            ></textarea>
        </div>
    </div>

    <div class="preview">
        <div class="email-frame">
            <div class="email-header">
                <h3>{emailPreview.subject}</h3>
                <span class="template-badge">{emailPreview.template}</span>
            </div>
            <div class="email-content">
                {emailPreview.content}
            </div>
        </div>
    </div>
</div>

<style>
    .email-preview {
        display: grid;
        grid-template-columns: 300px 1fr;
        gap: 2rem;
        padding: 2rem;
        max-width: 1200px;
        margin: 0 auto;
    }

    .sidebar {
        background: white;
        padding: 1.5rem;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }

    .render-info {
        color: #718096;
        font-size: 0.875rem;
        margin-bottom: 1.5rem;
    }

    .form-group {
        margin-bottom: 1.5rem;
    }

    label {
        display: block;
        margin-bottom: 0.5rem;
        color: #4a5568;
    }

    input, select, textarea {
        width: 100%;
        padding: 0.5rem;
        border: 1px solid #e2e8f0;
        border-radius: 4px;
    }

    .preview {
        background: #f7fafc;
        padding: 2rem;
        border-radius: 8px;
    }

    .email-frame {
        background: white;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        overflow: hidden;
    }

    .email-header {
        background: #f7fafc;
        padding: 1rem;
        border-bottom: 1px solid #e2e8f0;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .email-content {
        padding: 2rem;
        min-height: 300px;
    }

    .template-badge {
        background: #edf2f7;
        padding: 0.25rem 0.5rem;
        border-radius: 4px;
        font-size: 0.875rem;
        color: #4a5568;
    }

    h2, h3 {
        color: #2d3748;
        margin: 0 0 1rem 0;
    }
</style>