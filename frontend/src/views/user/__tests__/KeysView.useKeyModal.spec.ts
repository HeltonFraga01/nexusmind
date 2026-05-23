import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount, type VueWrapper } from '@vue/test-utils'
import { nextTick } from 'vue'

import KeysView from '../KeysView.vue'

const {
  list,
  reveal,
  update,
  create,
  deleteKey,
  toggleStatus,
  getPublicSettings,
  getDashboardApiKeysUsage,
  getAvailable,
  getUserGroupRates,
  showError,
  showWarning,
  showSuccess,
  showInfo,
  isCurrentStep,
  nextStep,
  copyToClipboard
} = vi.hoisted(() => ({
  list: vi.fn(),
  reveal: vi.fn(),
  update: vi.fn(),
  create: vi.fn(),
  deleteKey: vi.fn(),
  toggleStatus: vi.fn(),
  getPublicSettings: vi.fn(),
  getDashboardApiKeysUsage: vi.fn(),
  getAvailable: vi.fn(),
  getUserGroupRates: vi.fn(),
  showError: vi.fn(),
  showWarning: vi.fn(),
  showSuccess: vi.fn(),
  showInfo: vi.fn(),
  isCurrentStep: vi.fn(),
  nextStep: vi.fn(),
  copyToClipboard: vi.fn()
}))

vi.mock('@/api', () => ({
  keysAPI: {
    list,
    reveal,
    update,
    create,
    delete: deleteKey,
    toggleStatus
  },
  authAPI: {
    getPublicSettings
  },
  usageAPI: {
    getDashboardApiKeysUsage
  },
  userGroupsAPI: {
    getAvailable,
    getUserGroupRates
  }
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({ showError, showWarning, showSuccess, showInfo })
}))

vi.mock('@/stores/onboarding', () => ({
  useOnboardingStore: () => ({ isCurrentStep, nextStep })
}))

vi.mock('@/composables/useClipboard', () => ({
  useClipboard: () => ({ copyToClipboard })
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) =>
      ({
        'keys.useKey': 'Use Key'
      })[key] ?? key
  })
}))

const AppLayoutStub = { template: '<div><slot /></div>' }
const TablePageLayoutStub = {
  template: '<div><slot name="actions" /><slot name="filters" /><slot name="table" /></div>'
}
const DataTableStub = {
  props: ['data'],
  template: `
    <div>
      <div v-for="row in data" :key="row.id">
        <slot name="cell-actions" :row="row" />
      </div>
    </div>
  `
}
const EmptyStateStub = { template: '<div />' }
const PaginationStub = { template: '<div />' }
const BaseDialogStub = { template: '<div><slot /><slot name="footer" /></div>' }
const ConfirmDialogStub = { template: '<div />' }
const SelectStub = { template: '<div />' }
const SearchInputStub = { template: '<input />' }
const IconStub = { template: '<span />' }
const UseKeyModalStub = {
  props: ['show', 'apiKey', 'baseUrl', 'platform', 'allowMessagesDispatch'],
  template: '<div v-if="show" data-testid="use-key-modal">{{ platform || \'no-platform\' }}</div>'
}
const EndpointPopoverStub = { template: '<div />' }
const GroupBadgeStub = { template: '<div />' }
const GroupOptionItemStub = { template: '<div />' }

vi.mock('@/components/layout/AppLayout.vue', () => ({ default: AppLayoutStub }))
vi.mock('@/components/layout/TablePageLayout.vue', () => ({ default: TablePageLayoutStub }))
vi.mock('@/components/common/DataTable.vue', () => ({ default: DataTableStub }))
vi.mock('@/components/common/EmptyState.vue', () => ({ default: EmptyStateStub }))
vi.mock('@/components/common/Pagination.vue', () => ({ default: PaginationStub }))
vi.mock('@/components/common/BaseDialog.vue', () => ({ default: BaseDialogStub }))
vi.mock('@/components/common/ConfirmDialog.vue', () => ({ default: ConfirmDialogStub }))
vi.mock('@/components/common/Select.vue', () => ({ default: SelectStub }))
vi.mock('@/components/common/SearchInput.vue', () => ({ default: SearchInputStub }))
vi.mock('@/components/icons/Icon.vue', () => ({ default: IconStub }))
vi.mock('@/components/keys/UseKeyModal.vue', () => ({ default: UseKeyModalStub }))
vi.mock('@/components/keys/EndpointPopover.vue', () => ({ default: EndpointPopoverStub }))
vi.mock('@/components/common/GroupBadge.vue', () => ({ default: GroupBadgeStub }))
vi.mock('@/components/common/GroupOptionItem.vue', () => ({ default: GroupOptionItemStub }))

const group = {
  id: 7,
  name: 'OpenAI Group',
  description: 'OpenAI compatible group',
  platform: 'openai',
  rate_multiplier: 1,
  subscription_type: 'shared',
  sort_order: 0
}

const groupedKey = {
  id: 42,
  user_id: 1,
  key: 'sk-***test',
  name: 'Grouped key',
  group_id: group.id,
  group,
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
}

function mountView() {
  return mount(KeysView, {
    global: {
      stubs: {
        AppLayout: AppLayoutStub,
        TablePageLayout: TablePageLayoutStub,
        DataTable: DataTableStub,
        EmptyState: EmptyStateStub,
        Pagination: PaginationStub,
        BaseDialog: BaseDialogStub,
        ConfirmDialog: ConfirmDialogStub,
        Select: SelectStub,
        SearchInput: SearchInputStub,
        Icon: IconStub,
        UseKeyModal: UseKeyModalStub,
        EndpointPopover: EndpointPopoverStub,
        GroupBadge: GroupBadgeStub,
        GroupOptionItem: GroupOptionItemStub,
        Teleport: true
      }
    }
  })
}

describe('KeysView use key modal', () => {
  let wrapper: VueWrapper | undefined

  beforeEach(() => {
    list.mockReset()
    reveal.mockReset()
    update.mockReset()
    create.mockReset()
    deleteKey.mockReset()
    toggleStatus.mockReset()
    getPublicSettings.mockReset()
    getDashboardApiKeysUsage.mockReset()
    getAvailable.mockReset()
    getUserGroupRates.mockReset()
    showError.mockReset()
    showWarning.mockReset()
    showSuccess.mockReset()
    showInfo.mockReset()
    isCurrentStep.mockReset()
    nextStep.mockReset()
    copyToClipboard.mockReset()

    list.mockResolvedValue({
      items: [groupedKey],
      total: 1,
      page: 1,
      page_size: 10,
      pages: 1
    })
    reveal.mockResolvedValue({
      api_key: {
        ...groupedKey,
        group: undefined
      },
      raw_key: 'sk-revealed'
    })
    getDashboardApiKeysUsage.mockResolvedValue({ stats: {} })
    getAvailable.mockResolvedValue([group])
    getUserGroupRates.mockResolvedValue({})
    getPublicSettings.mockResolvedValue({
      api_base_url: 'https://example.com/v1',
      hide_ccs_import_button: true,
      site_name: 'NexusMind'
    })
    isCurrentStep.mockReturnValue(false)
  })

  afterEach(() => {
    wrapper?.unmount()
    wrapper = undefined
  })

  it('uses clicked row group when reveal omits embedded group data', async () => {
    wrapper = mountView()
    await flushPromises()

    const useKeyButton = wrapper.findAll('button').find((button) => button.text().includes('Use Key'))

    expect(useKeyButton).toBeDefined()
    await useKeyButton!.trigger('click')
    await flushPromises()
    await nextTick()

    expect(reveal).toHaveBeenCalledWith(groupedKey.id)
    expect(wrapper.get('[data-testid="use-key-modal"]').text()).toBe('openai')
  })
})
