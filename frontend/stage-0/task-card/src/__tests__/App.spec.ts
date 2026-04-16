import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import App from "../App.vue";

describe("TodoCard", () => {
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
    expect((checkbox.element as HTMLInputElement).checked).toBe(false);
    await checkbox.setValue(true);
    expect((checkbox.element as HTMLInputElement).checked).toBe(true);
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
