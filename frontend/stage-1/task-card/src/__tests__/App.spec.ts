import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import App from "../App.vue";

// ─── Stage 0: All original testids preserved ──────────────────────────────────
describe("TodoCard – Stage 0 (preserved)", () => {
  it("renders the card container", () => {
    const wrapper = mount(App);
    expect(wrapper.find('[data-testid="test-todo-card"]').exists()).toBe(true);
  });

  it("renders a non-empty title", () => {
    const wrapper = mount(App);
    expect(wrapper.find('[data-testid="test-todo-title"]').text().trim()).not.toBe("");
  });

  it("renders a non-empty description", () => {
    const wrapper = mount(App);
    expect(wrapper.find('[data-testid="test-todo-description"]').text().trim()).not.toBe("");
  });

  it("renders a priority badge", () => {
    const wrapper = mount(App);
    expect(wrapper.find('[data-testid="test-todo-priority"]').exists()).toBe(true);
  });

  it("renders a due date", () => {
    const wrapper = mount(App);
    expect(wrapper.find('[data-testid="test-todo-due-date"]').text().trim()).not.toBe("");
  });

  it("renders a time remaining value", () => {
    const wrapper = mount(App);
    expect(wrapper.find('[data-testid="test-todo-time-remaining"]').text().trim()).not.toBe("");
  });

  it("renders a status indicator", () => {
    const wrapper = mount(App);
    expect(wrapper.find('[data-testid="test-todo-status"]').exists()).toBe(true);
  });

  it("renders a checkbox that can be toggled", async () => {
    const wrapper = mount(App);
    const checkbox = wrapper.find<HTMLInputElement>('[data-testid="test-todo-complete-toggle"]');
    expect(checkbox.exists()).toBe(true);
    expect(checkbox.element.checked).toBe(false);
    await checkbox.setValue(true);
    expect(checkbox.element.checked).toBe(true);
  });

  it("renders the tags list with at least one tag", () => {
    const wrapper = mount(App);
    const tags = wrapper.find('[data-testid="test-todo-tags"]');
    expect(tags.exists()).toBe(true);
    expect(tags.findAll("li").length).toBeGreaterThan(0);
  });

  it("renders the edit button", () => {
    const wrapper = mount(App);
    expect(wrapper.find('[data-testid="test-todo-edit-button"]').exists()).toBe(true);
  });

  it("renders the delete button", () => {
    const wrapper = mount(App);
    expect(wrapper.find('[data-testid="test-todo-delete-button"]').exists()).toBe(true);
  });
});

// ─── Stage 1A: New features ───────────────────────────────────────────────────
describe("TodoCard – Stage 1A: Edit Mode", () => {
  it("shows edit form when Edit button is clicked", async () => {
    const wrapper = mount(App);
    await wrapper.find('[data-testid="test-todo-edit-button"]').trigger("click");
    expect(wrapper.find('[data-testid="test-todo-edit-form"]').exists()).toBe(true);
  });

  it("edit form contains all required inputs", async () => {
    const wrapper = mount(App);
    await wrapper.find('[data-testid="test-todo-edit-button"]').trigger("click");
    expect(wrapper.find('[data-testid="test-todo-edit-title-input"]').exists()).toBe(true);
    expect(wrapper.find('[data-testid="test-todo-edit-description-input"]').exists()).toBe(true);
    expect(wrapper.find('[data-testid="test-todo-edit-priority-select"]').exists()).toBe(true);
    expect(wrapper.find('[data-testid="test-todo-edit-due-date-input"]').exists()).toBe(true);
    expect(wrapper.find('[data-testid="test-todo-save-button"]').exists()).toBe(true);
    expect(wrapper.find('[data-testid="test-todo-cancel-button"]').exists()).toBe(true);
  });

  it("Save updates the card title", async () => {
    const wrapper = mount(App);
    await wrapper.find('[data-testid="test-todo-edit-button"]').trigger("click");
    await wrapper.find('[data-testid="test-todo-edit-title-input"]').setValue("New Title");
    await wrapper.find('[data-testid="test-todo-save-button"]').trigger("click");
    expect(wrapper.find('[data-testid="test-todo-title"]').text()).toBe("New Title");
  });

  it("Cancel restores previous title", async () => {
    const wrapper = mount(App);
    const originalTitle = wrapper.find('[data-testid="test-todo-title"]').text();
    await wrapper.find('[data-testid="test-todo-edit-button"]').trigger("click");
    await wrapper.find('[data-testid="test-todo-edit-title-input"]').setValue("Changed Title");
    await wrapper.find('[data-testid="test-todo-cancel-button"]').trigger("click");
    expect(wrapper.find('[data-testid="test-todo-title"]').text()).toBe(originalTitle);
  });

  it("card view is hidden while editing", async () => {
    const wrapper = mount(App);
    await wrapper.find('[data-testid="test-todo-edit-button"]').trigger("click");
    expect(wrapper.find('[data-testid="test-todo-title"]').exists()).toBe(false);
  });
});

describe("TodoCard – Stage 1A: Status Controls", () => {
  it("renders the status control", () => {
    const wrapper = mount(App);
    expect(wrapper.find('[data-testid="test-todo-status-control"]').exists()).toBe(true);
  });

  it("status control changes the status display", async () => {
    const wrapper = mount(App);
    await wrapper
      .find<HTMLSelectElement>('[data-testid="test-todo-status-control"]')
      .setValue("Done");
    expect(wrapper.find('[data-testid="test-todo-status"]').text()).toBe("Done");
  });

  it("checking the checkbox sets status to Done", async () => {
    const wrapper = mount(App);
    await wrapper.find('[data-testid="test-todo-complete-toggle"]').setValue(true);
    expect(wrapper.find('[data-testid="test-todo-status"]').text()).toBe("Done");
  });

  it("setting status to Done checks the checkbox", async () => {
    const wrapper = mount(App);
    await wrapper
      .find<HTMLSelectElement>('[data-testid="test-todo-status-control"]')
      .setValue("Done");
    const checkbox = wrapper.find<HTMLInputElement>('[data-testid="test-todo-complete-toggle"]');
    expect(checkbox.element.checked).toBe(true);
  });

  it("unchecking checkbox after Done reverts status to Pending", async () => {
    const wrapper = mount(App);
    await wrapper.find('[data-testid="test-todo-complete-toggle"]').setValue(true);
    await wrapper.find('[data-testid="test-todo-complete-toggle"]').setValue(false);
    expect(wrapper.find('[data-testid="test-todo-status"]').text()).toBe("Pending");
  });
});

describe("TodoCard – Stage 1A: Priority Indicator", () => {
  it("renders the priority indicator", () => {
    const wrapper = mount(App);
    expect(wrapper.find('[data-testid="test-todo-priority-indicator"]').exists()).toBe(true);
  });
});

describe("TodoCard – Stage 1A: Expand / Collapse", () => {
  it("renders the collapsible section", () => {
    const wrapper = mount(App);
    expect(wrapper.find('[data-testid="test-todo-collapsible-section"]').exists()).toBe(true);
  });

  it("renders the expand toggle when description is long", () => {
    const wrapper = mount(App);
    expect(wrapper.find('[data-testid="test-todo-expand-toggle"]').exists()).toBe(true);
  });

  it("expand toggle starts collapsed (aria-expanded=false)", () => {
    const wrapper = mount(App);
    expect(
      wrapper.find('[data-testid="test-todo-expand-toggle"]').attributes("aria-expanded"),
    ).toBe("false");
  });

  it("clicking expand toggle sets aria-expanded to true", async () => {
    const wrapper = mount(App);
    const toggle = wrapper.find('[data-testid="test-todo-expand-toggle"]');
    await toggle.trigger("click");
    expect(toggle.attributes("aria-expanded")).toBe("true");
  });

  it("clicking expand toggle again collapses back", async () => {
    const wrapper = mount(App);
    const toggle = wrapper.find('[data-testid="test-todo-expand-toggle"]');
    await toggle.trigger("click");
    await toggle.trigger("click");
    expect(toggle.attributes("aria-expanded")).toBe("false");
  });
});

describe("TodoCard – Stage 1A: Time Management", () => {
  it("time remaining shows 'Completed' when status is Done", async () => {
    const wrapper = mount(App);
    await wrapper
      .find<HTMLSelectElement>('[data-testid="test-todo-status-control"]')
      .setValue("Done");
    expect(wrapper.find('[data-testid="test-todo-time-remaining"]').text()).toBe("Completed");
  });
});
