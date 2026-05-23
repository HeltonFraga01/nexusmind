import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { defineComponent, nextTick } from 'vue'

import KeysView from '../KeysView.vue'

const {
  listMock,
  createMock,
  revealMock,
  updateMock,
  deleteMock,
  toggleStatusMock,
  getDashboardApiKeysUsageMock,
  getAvailableMock,
  getUserGroupRatesMock,
  getPublicSettingsMock,
  showError,
  showSuccess,
  onboardingIsCurrentStep,
  onboardingNextStep
} = vi.hoisted(() => ({
  listMock: vi.fn(),
  createMock: vi.fn(),
  revealMock: vi.fn(),
  updateMock: vi.fn(),
  deleteMock: vi.fn(),
  toggleStatusMock: vi.fn(),
  getDashboardApiKeysUsageMock: vi.fn(),
  getAvailableMock: vi.fn(),
  getUserGroupRatesMock: vi.fn(),
  getPublicSettingsMock: vi.fn(),
  showError: vi.fn(),
  showSuccess: vi.fn(),
  onboardingIsCurrentStep: vi.fn(),
  onboardingNextStep: vi.fn()
}))

vi.mock('@/api', () => ({
  keysAPI: {
    list: listMock,
    create: createMock,
    reveal: revealMock,
    update: updateMock,
    delete: deleteMock,
    toggleStatus: toggleStatusMock
  },
  authAPI: {
    getPublicSettings: getPublicSettingsMock
  },
  usageAPI: {
    getDashboardApiKeysUsage: getDashboardApiKeysUsageMock
  },
  userGroupsAPI: {
    getAvailable: getAvailableMock,
    getUserGroupRates: getUserGroupRatesMock
  }
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({ showError, showSuccess })
}))

vi.mock('@/stores/onboarding', () => ({
  useOnboardingStore: () => ({
    isCurrentStep: onboardingIsCurrentStep,
    nextStep: onboardingNextStep
  })
}))

vi.mock('@/composables/useClipboard', () => ({
  useClipboard: () => ({
    copyToClipboard: vi.fn().mockResolvedValue(true)
  })
}))

vi.mock('@/composables/usePersistedPageSize', () => ({
  getPersistedPageSize: () => 10
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key
    })
  }
})

const UseKeyModalStub = defineComponent({
  name: 'UseKeyModal',
  props: {
    show: { type: Boolean, default: false },
    platform: { type: String, default: null },
    allowMessagesDispatch: { type: Boolean, default: false }
  },
  template: `
    <div v-if="show" data-testid="use-key-modal">
      <span data-testid="use-key-platform">{{ platform || 'null' }}</span>
      <span data-testid="use-key-dispatch">{{ allowMessagesDispatch ? 'enabled' : 'disabled' }}</span>
      <p v-if="!platform">keys.useKeyModal.noGroupTitle</p>
    </div>
  `
})

const BaseDialogStub = defineComponent({
  name: 'BaseDialog',
  props: {
    show: { type: Boolean, default: false }
  },
  template: '<div v-if="show"><slot /></div>'
})

describe('KeysView', () => {
  beforeEach(() => {
    listMock.mockReset()
    createMock.mockReset()
    revealMock.mockReset()
    updateMock.mockReset()
    deleteMock.mockReset()
    toggleStatusMock.mockReset()
    getDashboardApiKeysUsageMock.mockReset()
    getAvailableMock.mockReset()
    getUserGroupRatesMock.mockReset()
    getPublicSettingsMock.mockReset()
    showError.mockReset()
    showSuccess.mockReset()
    onboardingIsCurrentStep.mockReset()
    onboardingNextStep.mockReset()

    listMock.mockResolvedValue({ items: [], total: 0, pages: 0 })
    getAvailableMock.mockResolvedValue([
      {
        id: 7,
        name: 'Anthropic Group',
        description: null,
        platform: 'anthropic',
        rate_multiplier: 1,
        is_exclusive: false,
        status: 'active',
        subscription_type: 'free',
        daily_limit_usd: null,
        weekly_limit_usd: null,
        monthly_limit_usd: null,
        allow_image_generation: false,
        image_rate_independent: false,
        image_rate_multiplier: 1,
        image_price_1k: null,
        image_price_2k: null,
        image_price_4k: null,
        claude_code_only: false,
        fallback_group_id: null,
        fallback_group_id_on_invalid_request: null,
        allow_messages_dispatch: true,
        require_oauth_only: false,
        require_privacy_set: false,
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z'
      }
    ])
    getUserGroupRatesMock.mockResolvedValue({ 7: 1 })
    getPublicSettingsMock.mockResolvedValue({
      site_name: 'NexusMind',
      api_base_url: 'https://nexusmind.digital/v1'
    })
    getDashboardApiKeysUsageMock.mockResolvedValue({ stats: {} })
    onboardingIsCurrentStep.mockReturnValue(false)

    createMock.mockResolvedValue({
      api_key: {
        id: 101,
        user_id: 1,
        key: 'sk-****',
        name: 'new key',
        group_id: 7,
        status: 'active',
        ip_whitelist: [],
        ip_blacklist: [],
        last_used_at: null,
        quota: 0,
        quota_used: 0,
        expires_at: null,
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
        rate_limit_5h: 0,
        rate_limit_1d: 0,
        rate_limit_7d: 0,
        usage_5h: 0,
        usage_1d: 0,
        usage_7d: 0,
        window_5h_start: null,
        window_1d_start: null,
        window_7d_start: null,
        reset_5h_at: null,
        reset_1d_at: null,
        reset_7d_at: null
      },
      raw_key: 'sk-live-created'
    })
  })

  it('keeps group platform available when opening Use API Key modal right after create', async () => {
    const wrapper = mount(KeysView, {
      global: {
        stubs: {
          AppLayout: { template: '<div><slot /></div>' },
          TablePageLayout: { template: '<div><slot name="filters" /><slot name="actions" /><slot name="table" /><slot name="pagination" /></div>' },
          DataTable: { template: '<div />' },
          Pagination: { template: '<div />' },
          BaseDialog: BaseDialogStub,
          ConfirmDialog: { template: '<div />' },
          EmptyState: { template: '<div />' },
          Select: { template: '<div />' },
          SearchInput: { template: '<div />' },
          Icon: { template: '<span />' },
          EndpointPopover: { template: '<div />' },
          GroupBadge: { template: '<div />' },
          GroupOptionItem: { template: '<div />' },
          UseKeyModal: UseKeyModalStub
        }
      }
    })

    await flushPromises()

    ;(wrapper.vm as any).showCreateModal = true
    ;(wrapper.vm as any).formData.name = 'new key'
    ;(wrapper.vm as any).formData.group_id = 7

    await nextTick()
    await (wrapper.vm as any).handleSubmit()
    await flushPromises()

    expect(createMock).toHaveBeenCalledOnce()
    expect((wrapper.vm as any).showCreateModal).toBe(false)
    expect((wrapper.vm as any).showUseKeyModal).toBe(true)

    const modal = wrapper.get('[data-testid="use-key-modal"]')
    expect(modal.exists()).toBe(true)
    expect(wrapper.get('[data-testid="use-key-platform"]').text()).toBe('anthropic')
    expect(wrapper.get('[data-testid="use-key-dispatch"]').text()).toBe('enabled')
    expect(modal.text()).not.toContain('keys.useKeyModal.noGroupTitle')
  })
})
