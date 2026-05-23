<template>
  <!-- Admin custom content mode -->
  <div v-if="homeContent" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Premium NexusMind State-of-the-Art Landing Page -->
  <div v-else class="nm-home">
    <div class="nm-bg-grid"></div>
    <div class="nm-bg-glow"></div>

    <!-- ── Header ─────────────────────────────────────────── -->
    <header class="nm-header">
      <div class="nm-wrap nm-header-inner">
        <a href="/home" class="nm-brand">
          <div class="nm-brand-logo">
            <img :src="siteLogo || '/logo.png'" :alt="siteName" />
          </div>
          <div class="nm-brand-name">
            <b>{{ siteName }}</b>
            <span>AI ROUTING · BR</span>
          </div>
        </a>

        <nav class="nm-nav">
          <a href="#nm-how" class="nm-nav-link">$ {{ c.navHow }}</a>
          <a href="#nm-math" class="nm-nav-link">$ {{ c.navMath }}</a>
          <a href="#nm-pricing" class="nm-nav-link">$ {{ c.navPricing }}</a>
          <a href="#nm-faq" class="nm-nav-link">$ {{ c.navFaq }}</a>
        </nav>

        <div class="nm-actions">
          <LocaleSwitcher />
          <router-link v-if="isAuthenticated" :to="dashboardPath" class="nm-btn nm-btn-ghost">
            {{ c.btnDashboard }}
          </router-link>
          <router-link v-else to="/login" class="nm-btn nm-btn-ghost">
            {{ c.btnLogin }}
          </router-link>
          <a href="#nm-pricing" class="nm-btn nm-btn-primary">
            {{ c.btnStart }} <span class="nm-arr">→</span>
          </a>
        </div>
      </div>
    </header>

    <!-- ── Main ───────────────────────────────────────────── -->
    <main>

      <!-- 🌌 Hero Section -->
      <section id="nm-top" class="nm-hero nm-wrap">
        <div class="nm-hero-grid">
          <div>
            <div class="nm-eyebrow">
              <span class="nm-eyebrow-dot animate-pulse"></span>
              <span>{{ siteName }} · <span class="nm-eyebrow-v">v3.0 SEAS</span> · BR</span>
            </div>
            <h1 class="nm-h1" v-html="c.heroTitle"></h1>
            <p class="nm-lede">{{ c.heroLede }}</p>
            
            <div class="nm-cta-row">
              <router-link :to="isAuthenticated ? dashboardPath : '/login'" class="nm-btn nm-btn-primary nm-btn-lg">
                {{ isAuthenticated ? c.ctaOpenDashboard : c.ctaGetStarted }}
                <span class="nm-arr">→</span>
              </router-link>
              <a href="#nm-how" class="nm-btn nm-btn-ghost nm-btn-lg">{{ c.ctaHowItWorks }}</a>
            </div>
            
            <div class="nm-pill-row">
              <span class="nm-pill-item">
                <span class="nm-dot-green"></span> {{ c.heroPill1 }}
              </span>
              <span class="nm-pill-item">
                <span class="nm-ok">●</span> {{ c.heroPill2 }}
              </span>
              <span class="nm-pill-item">{{ c.heroPill3 }}</span>
            </div>
          </div>

          <!-- 🖥️ Routing Console (Real-time Simulation) -->
          <div class="nm-console">
            <div class="nm-console-head">
              <span class="nm-console-dots"><i class="red" /><i class="yellow" /><i class="green" /></span>
              <span class="nm-console-title">gateway.nexusmind.digital</span>
              <span class="nm-console-live">
                <i class="nm-live-dot"></i> {{ c.consoleLive }}
              </span>
            </div>
            
            <div class="nm-console-body">
              <!-- SVG diagram (packets + active lanes) -->
              <svg class="nm-diagram" viewBox="0 0 560 320" preserveAspectRatio="xMidYMid meet">
                <defs>
                  <linearGradient id="nmG1" x1="0" y1="0" x2="1" y2="0">
                    <stop offset="0%" stop-color="#2dd4bf" />
                    <stop offset="100%" stop-color="#7c3aed" />
                  </linearGradient>
                  <radialGradient id="nmCoreGlow" cx="50%" cy="50%" r="50%">
                    <stop offset="0%" stop-color="rgba(45,212,191,0.4)" />
                    <stop offset="100%" stop-color="rgba(45,212,191,0)" />
                  </radialGradient>
                </defs>
                <g opacity="0.06">
                  <line v-for="n in 7" :key="n" x1="0" :y1="n*46+20" x2="560" :y2="n*46+20" stroke="#2dd4bf" stroke-dasharray="2 6" />
                </g>
                <text x="14" y="22" class="nm-svg-dim">$ ingress</text>
                <text x="540" y="22" class="nm-svg-dim" text-anchor="end">$ providers</text>

                <!-- Request nodes -->
                <g v-for="(req, i) in nmRequests" :key="i" :opacity="activeRoute === i ? 1 : 0.4">
                  <rect :x="14" :y="req.y-12" width="180" height="24" rx="4" fill="#03050a" :stroke="activeRoute === i ? '#2dd4bf' : 'rgba(255,255,255,0.08)'" />
                  <circle cx="26" :cy="req.y" r="3" :fill="activeRoute === i ? '#2dd4bf' : '#7c3aed'" />
                  <text x="36" :y="req.y+4" :class="['nm-svg-label', activeRoute === i ? 'nm-svg-teal' : 'nm-svg-dim']">{{ req.label }}</text>
                </g>

                <!-- Lanes to core -->
                <path
                  v-for="(req, i) in nmRequests"
                  :key="'l'+i"
                  :d="`M194 ${req.y} C 240 ${req.y}, 240 160, 280 160`"
                  stroke="url(#nmG1)"
                  :class="['nm-lane', activeRoute === i ? 'nm-lane-active' : '']"
                />

                <!-- Core (router) -->
                <circle cx="280" cy="160" r="58" fill="url(#nmCoreGlow)" />
                <rect x="235" y="136" width="90" height="48" rx="8" fill="#080c14" stroke="#2dd4bf" stroke-width="1.5" />
                <text x="280" y="154" class="nm-svg-teal" text-anchor="middle" font-size="11" font-weight="700" font-family="JetBrains Mono,monospace" style="letter-spacing: 0.08em">NEXUS ENGINE</text>
                <text x="280" y="170" class="nm-svg-dim" text-anchor="middle" font-size="9">nexusmind.digital</text>

                <!-- Provider lanes + boxes -->
                <g v-for="(prov, i) in nmProviders" :key="'p'+i">
                  <path
                    :d="`M320 160 C 380 160, 380 ${prov.y}, 420 ${prov.y}`"
                    stroke="url(#nmG1)"
                    :class="['nm-lane', activeRoute === i ? 'nm-lane-active' : '']"
                  />
                  <rect :x="420" :y="prov.y-12" width="126" height="24" rx="4" :class="['nm-prov-svg', activeRoute === i ? 'nm-prov-svg-active' : '']" />
                  <rect :x="426" :y="prov.y-7" width="14" height="14" rx="2" :fill="activeRoute === i ? '#2dd4bf' : '#111827'" />
                  <text :x="448" :y="prov.y+4" :class="['nm-svg-label', activeRoute === i ? 'nm-svg-teal' : 'nm-svg-dim']">{{ prov.name }}</text>
                  <text x="540" :y="prov.y+4" class="nm-svg-dim" text-anchor="end" font-size="9">{{ activeRoute === i ? '200 OK' : '—' }}</text>
                </g>

                <!-- Traveling packet -->
                <circle v-if="activeRoute === 0" :key="'pk0-'+packetTick" class="nm-packet nm-packet-1" r="4.5" />
                <circle v-if="activeRoute === 1" :key="'pk1-'+packetTick" class="nm-packet nm-packet-2" r="4.5" />
                <circle v-if="activeRoute === 2" :key="'pk2-'+packetTick" class="nm-packet nm-packet-3" r="4.5" />
              </svg>

              <!-- Routing target tag -->
              <div class="nm-routing-target">
                <span class="nm-rt-label">$ {{ c.consoleRouting }}</span>
                <span class="nm-rt-val">{{ nmRequests[activeRoute].method }}</span>
                <span class="nm-rt-label">→</span>
                <span class="nm-rt-val">{{ nmRequests[activeRoute].prov }}</span>
                <span class="nm-rt-val">{{ nmRequests[activeRoute].model }}</span>
              </div>

              <!-- Live log feed -->
              <div class="nm-log-feed">
                <div v-for="(log, i) in logs" :key="log.id" class="nm-log-line" :style="{ opacity: 1 - i * 0.32 }">
                  <span class="nm-log-t">[{{ log.t }}]</span>
                  <span class="nm-log-method">{{ log.path }}</span>
                  <span class="nm-log-arr">→</span>
                  <span class="nm-log-prov">{{ log.prov }}</span>
                  <span class="nm-log-model">{{ log.model }}</span>
                  <span class="nm-log-latency">{{ log.latency }}ms</span>
                  <span class="nm-log-ok">200</span>
                </div>
              </div>

              <!-- Meta specs -->
              <div class="nm-console-meta">
                <div class="nm-meta-cell">
                  <div class="nm-meta-k">{{ c.metaReqs }}</div>
                  <div class="nm-meta-v">428<span class="nm-meta-u">/s</span></div>
                </div>
                <div class="nm-meta-cell">
                  <div class="nm-meta-k">{{ c.metaLatency }}</div>
                  <div class="nm-meta-v">312<span class="nm-meta-u">ms</span></div>
                </div>
                <div class="nm-meta-cell">
                  <div class="nm-meta-k">{{ c.metaSavings }}</div>
                  <div class="nm-meta-v">R$ 851<span class="nm-meta-u">/{{ c.metaMonth }}</span></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 🎥 Tutorial Section -->
      <section id="nm-tutorial" class="nm-section nm-wrap border-t">
        <div class="nm-s-label"><span class="nm-sq"></span> $ tutorial <span class="nm-tic">--duration=90s</span></div>
        <h2 class="nm-s-title">{{ c.tutTitle }}</h2>
        <p class="nm-s-sub">{{ c.tutSubtitle }}</p>

        <div class="nm-tut-grid">
          <div>
            <div class="nm-video-wrap">
              <iframe
                v-if="videoEmbedUrl"
                :src="videoEmbedUrl"
                title="Tutorial NexusMind"
                allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                allowfullscreen
              ></iframe>
              <div v-else class="nm-video-placeholder">
                <div class="nm-video-badge">✦ TUTORIAL · 90s</div>
                <div class="nm-video-time">{{ c.tutNoAudio }}</div>
                <button class="nm-play-btn" aria-label="Reproduzir tutorial">
                  <svg viewBox="0 0 24 24" width="32" height="32" fill="currentColor">
                    <path d="M8 5v14l11-7z" />
                  </svg>
                </button>
                <div class="nm-video-hint" v-html="c.tutHint"></div>
                <div class="nm-video-strip">
                  <span class="nm-vs-left"><span class="nm-vs-dot"></span> {{ siteName.toUpperCase() }} · /tutorial</span>
                  <span class="nm-vs-progress"></span>
                  <span>00:22 / 01:30</span>
                </div>
              </div>
            </div>

            <!-- 🎁 Free Trial Action Banner (Positioned right near/under the video player) -->
            <div class="nm-action-banner nm-reveal">
              <div class="nm-ab-content">
                <div class="nm-ab-icon">🎁</div>
                <div class="nm-ab-text">
                  <h4>{{ c.freeTrialTitle }}</h4>
                  <p>{{ c.freeTrialDesc }}</p>
                </div>
              </div>
              <router-link :to="isAuthenticated ? dashboardPath : '/login'" class="nm-btn nm-btn-primary nm-btn-lg">
                {{ c.freeTrialBtn }} <span class="nm-arr">→</span>
              </router-link>
            </div>
          </div>

          <div class="nm-tut-chapters">
            <h4 class="nm-tut-h">$ chapters --list</h4>
            <ul class="nm-tut-list">
              <li v-for="(ch, i) in translatedChapters" :key="i">
                <span class="nm-tut-time">{{ ch.t }}</span>
                <span>{{ ch.txt }}</span>
              </li>
            </ul>
          </div>
        </div>
      </section>

      <!-- ⚡ Como Funciona Section -->
      <section id="nm-how" class="nm-section nm-wrap border-t">
        <div class="nm-s-label"><span class="nm-sq"></span> $ como_funciona <span class="nm-tic">--steps=3</span></div>
        <h2 class="nm-s-title" v-html="c.howTitle"></h2>
        <p class="nm-s-sub">{{ c.howSubtitle }}</p>

        <div class="nm-steps">
          <article class="nm-step">
            <div class="nm-step-num"><b>01</b> /// {{ c.step1Num }}</div>
            <h3>{{ c.step1Title }}</h3>
            <p>{{ c.step1Desc }}</p>
            <div class="nm-codeblock">
              <span class="nm-c"># panel.nexusmind.digital</span><br />
              <span class="nm-k">sk_live</span>_<span class="nm-v">7c9f...8b2a</span>
            </div>
          </article>
          
          <article class="nm-step">
            <div class="nm-step-num"><b>02</b> /// {{ c.step2Num }}</div>
            <h3>{{ c.step2Title }}</h3>
            <p>{{ c.step2Desc }}</p>
            <div class="nm-codeblock">
              <span class="nm-k">ANTHROPIC_BASE_URL</span>=<br />
              <span class="nm-s">"https://gateway.nexusmind.digital"</span>
            </div>
          </article>
          
          <article class="nm-step">
            <div class="nm-step-num"><b>03</b> /// {{ c.step3Num }}</div>
            <h3>{{ c.step3Title }}</h3>
            <p>{{ c.step3Desc }}</p>
            <div class="nm-codeblock">
              <span class="nm-c">→ model=claude-opus-4.7</span><br />
              <span class="nm-c">→ routed to provider AN_03</span><br />
              <span class="nm-s">← 200 OK · 482ms</span>
            </div>
          </article>
        </div>
      </section>

      <!-- 📊 Matemática Section (FinOps Comparison) -->
      <section id="nm-math" class="nm-section nm-wrap border-t">
        <div class="nm-s-label"><span class="nm-sq"></span> $ matematica <span class="nm-tic">--why=cheaper</span></div>
        <h2 class="nm-s-title" v-html="c.mathTitle"></h2>
        <p class="nm-s-sub">{{ c.mathSubtitle }}</p>

        <div class="nm-math-grid">
          <div class="nm-compare">
            <div class="nm-cmp-row nm-cmp-header">
              <div class="nm-cmp-cell">{{ c.mathCol1 }}</div>
              <div class="nm-cmp-cell">{{ c.mathCol2 }}</div>
              <div class="nm-cmp-cell nm-cmp-ours">{{ c.mathCol3 }} {{ siteName }}</div>
            </div>
            <div v-for="row in rcMathRows" :key="row.label" class="nm-cmp-row">
              <div class="nm-cmp-cell nm-cmp-label">{{ row.label }}</div>
              <div class="nm-cmp-cell"><span class="nm-strike">{{ row.direct }}</span></div>
              <div class="nm-cmp-cell nm-cmp-ours">{{ c.mathIncluded }}</div>
            </div>
            <div class="nm-cmp-row nm-cmp-total">
              <div class="nm-cmp-cell">{{ c.mathColTotal }}</div>
              <div class="nm-cmp-cell"><span class="nm-strike">R$ 1.020</span></div>
              <div class="nm-cmp-cell nm-cmp-ours">R$ 249</div>
            </div>
          </div>

          <div class="nm-math-points">
            <div
              v-for="(pt, i) in translatedMathPoints"
              :key="i"
              class="nm-math-pt"
              :class="{ 'nm-math-pt-active': activeMath === i }"
              @mouseenter="activeMath = i"
            >
              <div class="nm-pt-tag">{{ pt.tag }}</div>
              <h4>{{ pt.title }}</h4>
              <p>{{ pt.body }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- 🔌 Integrações Section -->
      <section id="nm-integrations" class="nm-section nm-wrap border-t">
        <div class="nm-s-label"><span class="nm-sq"></span> $ integracoes <span class="nm-tic">--openai-compat=true</span></div>
        <h2 class="nm-s-title" v-html="c.intTitle"></h2>
        <p class="nm-s-sub">{{ c.intSubtitle }}</p>

        <div class="nm-int-grid">
          <article
            v-for="(t, i) in translatedIntegrations"
            :key="i"
            class="nm-int-card"
          >
            <span v-if="t.tested" class="nm-int-tested"><span class="nm-int-d animate-pulse"></span> {{ c.intTested }}</span>
            <div class="nm-int-row">
              <div class="nm-int-icon">{{ t.icon }}</div>
              <div class="nm-int-name">{{ t.name }}</div>
            </div>
            <div class="nm-int-where">{{ t.where }}</div>
            <div class="nm-int-env"><b>{{ t.env }}</b>=<span class="nm-int-url">gateway.nexusmind.digital</span></div>
          </article>
        </div>

        <div class="nm-int-more">
          <b>+ {{ c.intMoreTitle }}:</b> {{ c.intMoreText }}
        </div>
      </section>

      <!-- 🛡️ Providers Section -->
      <section id="nm-providers" class="nm-section nm-wrap border-t">
        <div class="nm-s-label"><span class="nm-sq"></span> $ providers <span class="nm-tic">--catalog --all</span></div>
        <h2 class="nm-s-title" v-html="c.provTitle"></h2>
        <p class="nm-s-sub">{{ c.provSubtitle }}</p>
        
        <div class="nm-providers">
          <article v-for="pv in translatedProviderCards" :key="pv.code" class="nm-prov-card">
            <span class="nm-prov-status"><span class="nm-prov-dot animate-ping-slow"></span> ON</span>
            <div class="nm-prov-ico" :style="{ color: pv.color, borderColor: pv.color + '33' }">{{ pv.code }}</div>
            <h4>{{ pv.name }}</h4>
            <div class="nm-prov-models">
              <b>{{ pv.family }}</b><br />
              <div class="nm-prov-pills">
                <span v-for="m in pv.modelList" :key="m">{{ m }}</span>
              </div>
            </div>
          </article>
        </div>

        <div class="nm-providers-footer">
          <span v-html="c.provFooterText"></span>
        </div>
      </section>

      <!-- 💰 Pricing Section -->
      <section id="nm-pricing" class="nm-section nm-wrap border-t">
        <div class="nm-s-label"><span class="nm-sq"></span> $ pricing <span class="nm-tic">--currency=BRL</span></div>
        <h2 class="nm-s-title" v-html="c.priceTitle"></h2>
        <p class="nm-s-sub">{{ c.priceSubtitle }}</p>
        
        <div class="nm-pricing-grid">
          <article
            v-for="plan in translatedPlans"
            :key="plan.name"
            class="nm-plan"
            :class="{ 'nm-plan-feat': plan.feat }"
          >
            <span v-if="plan.feat" class="nm-plan-badge">{{ c.planFeatured }}</span>
            <div class="nm-plan-type">{{ plan.type }}</div>
            <h3>{{ plan.name }}</h3>
            <p class="nm-plan-desc">{{ plan.desc }}</p>
            <div class="nm-plan-price">
              <span class="nm-price-cur">R$</span>
              <span class="nm-price-amt">{{ plan.price }}</span>
              <span class="nm-price-per">/{{ c.metaMonth }}</span>
            </div>
            <div v-if="plan.oldPrice" class="nm-plan-old">{{ c.planEquivalent }} · R$ {{ plan.oldPrice }}/{{ c.metaMonth }}</div>
            <div v-else class="nm-plan-old nm-plan-old-spacer">·</div>
            
            <ul class="nm-plan-features">
              <li v-for="f in plan.features" :key="f">{{ f }}</li>
              <li v-for="f in plan.muted" :key="f" class="nm-feat-muted">{{ f }}</li>
            </ul>
            
            <router-link
              to="/purchase?tab=subscription"
              class="nm-btn nm-btn-lg"
              :class="plan.feat ? 'nm-btn-primary' : 'nm-btn-dark'"
              style="justify-content: center; display: flex;"
            >
              {{ plan.feat ? c.planBtnFeat : c.planBtnNormal + ' ' + plan.name }}
              <span class="nm-arr">→</span>
            </router-link>
          </article>
        </div>
      </section>

      <!-- 👥 Testimonials (Proof) Section -->
      <section id="nm-proof" class="nm-section nm-wrap border-t">
        <div class="nm-s-label"><span class="nm-sq"></span> $ devs --aprovam <span class="nm-tic">--count=100+</span></div>
        <h2 class="nm-s-title" v-html="c.proofTitle"></h2>
        <p class="nm-s-sub">{{ c.proofSubtitle }}</p>
        
        <div class="nm-proof-grid">
          <article v-for="(tm, i) in translatedTestimonials" :key="i" class="nm-tcard">
            <div class="nm-tcard-meta">
              <span>{{ tm.meta[0] }}</span>
              <span class="nm-t-sep">/</span>
              <span class="nm-t-ec">{{ tm.meta[1] }}</span>
            </div>
            <p class="nm-tcard-quote">"{{ tm.quote }}"</p>
            <div class="nm-tcard-who">
              <div class="nm-tcard-av">{{ tm.initial }}</div>
              <div>
                <b>{{ tm.name }}</b>
                <span>{{ tm.role }}</span>
              </div>
            </div>
          </article>
        </div>
      </section>

      <!-- 💬 FAQ Section -->
      <section id="nm-faq" class="nm-section nm-wrap border-t">
        <div class="nm-s-label"><span class="nm-sq"></span> $ faq <span class="nm-tic">--top=8</span></div>
        <h2 class="nm-s-title" v-html="c.faqTitle"></h2>
        <p class="nm-s-sub">{{ c.faqSubtitle }}</p>
        
        <div class="nm-faq-grid">
          <div
            v-for="(faq, i) in translatedFaqs"
            :key="i"
            class="nm-faq-item"
            :class="{ 'nm-faq-open': openFaq === i }"
          >
            <button @click="toggleFaq(i)">
              <span>
                <span class="nm-faq-num">Q.{{ String(i + 1).padStart(2, '0') }}</span>
                {{ faq.q }}
              </span>
              <span class="nm-faq-plus">+</span>
            </button>
            <div class="nm-faq-answer"><p>{{ faq.a }}</p></div>
          </div>
        </div>
      </section>

      <!-- 🚀 Final CTA Section -->
      <section class="nm-section nm-wrap border-t">
        <div class="nm-final-cta">
          <div class="nm-s-label" style="justify-content: center;">
            <span class="nm-sq"></span> $ proximo_passo <span class="nm-tic">--ready</span>
          </div>
          <h2 class="nm-s-title">{{ c.finalTitle }}</h2>
          <p class="nm-cta-p">{{ c.finalSubtitle }}</p>
          <div class="nm-cta-btns">
            <router-link to="/purchase?tab=subscription" class="nm-btn nm-btn-primary nm-btn-lg">
              {{ c.ctaGetStarted }} <span class="nm-arr">→</span>
            </router-link>
            <a :href="contactUrl" class="nm-btn nm-btn-ghost nm-btn-lg">{{ c.ctaWhatsApp }}</a>
          </div>
        </div>
      </section>

    </main>

    <!-- ── Footer ──────────────────────────────────────────── -->
    <footer class="nm-footer">
      <div class="nm-wrap nm-footer-grid">
        <div class="nm-footer-col nm-footer-about">
          <div class="nm-brand">
            <div class="nm-brand-logo">
              <img :src="siteLogo || '/logo.png'" alt="" />
            </div>
            <div class="nm-brand-name">
              <b>{{ siteName }}</b>
              <span>AI ROUTING · BR</span>
            </div>
          </div>
          <p class="nm-footer-desc">{{ c.footerDesc }}</p>
        </div>
        
        <div class="nm-footer-col">
          <h5>{{ c.footerProduct }}</h5>
          <a href="#nm-how">{{ c.navHow }}</a>
          <a href="#nm-providers">$ providers</a>
          <a href="#nm-pricing">{{ c.navPricing }}</a>
          <router-link :to="dashboardPath">{{ c.btnDashboard }}</router-link>
        </div>
        
        <div class="nm-footer-col">
          <h5>{{ c.footerChannels }}</h5>
          <a :href="contactUrl" target="_blank" rel="noopener noreferrer">WhatsApp</a>
          <a href="https://youtube.com" target="_blank" rel="noopener noreferrer">YouTube</a>
          <a href="https://instagram.com" target="_blank" rel="noopener noreferrer">Instagram</a>
        </div>
        
        <div class="nm-footer-col">
          <h5>Legal</h5>
          <a href="#">{{ c.legalTerms }}</a>
          <a href="#">{{ c.legalPrivacy }}</a>
          <a href="#">Status</a>
          <a href="#">{{ c.legalRefunds }}</a>
        </div>
      </div>
      
      <div class="nm-wrap nm-footer-bot">
        <span>© {{ currentYear }} {{ siteName }} · {{ c.footerBottomText }}</span>
        <span>nexusmind.digital</span>
      </div>
    </footer>

    <!-- ── Floating CTA Sticky Pill ── -->
    <div v-if="!floatDismissed" class="nm-float-cta" :class="{ 'nm-float-visible': floatVisible }">
      <div class="nm-float-inner">
        <span class="nm-float-prompt">{{ '>_' }} {{ siteName }}</span>
        <span class="nm-float-sep">·</span>
        <a :href="contactUrl" target="_blank" rel="noopener noreferrer" class="nm-float-link">WhatsApp</a>
        <a href="https://youtube.com" class="nm-float-link" target="_blank">YouTube</a>
        <a href="https://instagram.com" class="nm-float-link" target="_blank">Instagram</a>
        <router-link to="/purchase?tab=subscription" class="nm-float-btn">{{ c.floatBtn }} <span>→</span></router-link>
        <button class="nm-float-close" @click="floatDismissed = true" aria-label="Dispensar">×</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import { useAuthStore, useAppStore } from '@/stores'

const { locale } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()

// Computed locale check supporting pt-BR, zh, en
const isPt = computed(() => locale.value.startsWith('pt'))
const isZh = computed(() => locale.value.startsWith('zh'))

// ── Auth & Settings ────────────────────────────────────────
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => (isAdmin.value ? '/admin/dashboard' : '/dashboard'))
const currentYear = computed(() => new Date().getFullYear())

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'NexusMind')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const contactUrl = computed(() => appStore.cachedPublicSettings?.contact_info || 'https://wa.me/')

const isHomeContentUrl = computed(() => {
  const c = homeContent.value.trim()
  return c.startsWith('http://') || c.startsWith('https://')
})

// ── Computed Translations Dictionary ───────────────────────
const c = computed(() => {
  if (isPt.value) {
    return {
      navHow: "como_funciona",
      navMath: "matemática",
      navPricing: "preço",
      navFaq: "faq",
      btnDashboard: "Painel",
      btnLogin: "Entrar",
      btnStart: "Começar",
      heroTitle: "Uma chave.<br /><span class='nm-accent'>Todos os modelos</span><br />de IA.",
      heroLede: "Roteamento inteligente para Claude, GPT e Gemini — no real, com Pix. Painel pra monitorar tokens, custo e failover automático quando um provider engasga.",
      ctaGetStarted: "Quero minha chave",
      ctaOpenDashboard: "Abrir painel",
      ctaHowItWorks: "Como funciona",
      heroPill1: "3 provedores · 14 modelos",
      heroPill2: "99.94% uptime · 30d",
      heroPill3: "★ 4.9 · 100+ avaliações",
      consoleLive: "ROTEAMENTO",
      consoleRouting: "roteando agora",
      metaReqs: "req/s",
      metaLatency: "latência p95",
      metaSavings: "economia média",
      metaMonth: "mês",
      tutTitle: "Veja em 90 segundos.",
      tutSubtitle: "Do Pix até a primeira chamada. Sem áudio, sem firula — só os cliques na tela.",
      tutNoAudio: "SEM ÁUDIO",
      tutHint: "Cole o link do vídeo em <code>videoEmbedUrl</code> no código.",
      howTitle: "Em 3 passos.<br />Sem firula.",
      howSubtitle: "Você não precisa migrar nada, criar conta em provedor gringo, nem pagar dólar. Só trocar a base URL.",
      step1Num: "CRIAR CHAVE",
      step1Title: "Pagar e receber a chave",
      step1Desc: "Pix confirma na hora. A chave aparece no seu painel, no mesmo e-mail da compra.",
      step2Num: "CONFIGURAR",
      step2Title: "Trocar a base URL",
      step2Desc: "Cole a chave e a base URL no seu Cursor, Claude Code, script ou SDK. Funciona em qualquer cliente compatível com OpenAI.",
      step3Num: "USAR",
      step3Title: "Pedir e roteamos",
      step3Desc: "Seu pedido chega no roteador, escolhemos o provedor saudável mais barato para o modelo solicitado e devolvemos a resposta.",
      mathTitle: "Por que custa <span class='nm-accent'>R$ 249</span><br />no lugar de R$ 1.020?",
      mathSubtitle: "Sem pegadinha. Mesmos modelos, mesma API, mesma performance. A diferença vem de três fatores.",
      mathCol1: "Assinatura mensal",
      mathCol2: "Direto no provedor",
      mathCol3: "Via",
      mathColTotal: "Total / mês",
      mathIncluded: "incluso",
      intTitle: "Cole a chave onde<br />você já trabalha.",
      intSubtitle: "Trocou a base URL e a chave, funcionou. Qualquer cliente compatível com OpenAI funciona — sem migração, sem alterações complexas.",
      intTested: "Testado",
      intMoreTitle: "também roda em",
      intMoreText: "Zed · Vercel AI SDK · LlamaIndex · clientes OpenRouter · Goose · Roo Code · Cody · scripts Python/Node — qualquer ferramenta que aceite uma chave compatível com OpenAI.",
      provTitle: "Mais de 100 modelos<br />na mesma chave.",
      provSubtitle: "A homepage não precisa listar tudo: o ponto é que o FragaAi roteia entre famílias premium, coding, reasoning, multimodal e open-source via uma base OpenAI/Anthropic compatível.",
      provFooterText: "<b>Catálogo dinâmico:</b> além dos modelos acima, o gateway expõe combos como balanced-load, high-availability, coding, reasoning, economy, opus, flash e raiko. Novos modelos entram no roteador sem você trocar integração.",
      priceTitle: "Três planos.<br />Zero letra miúda.",
      priceSubtitle: "Cancele quando quiser. Pagamento facilitado via Pix ou cartão.",
      planFeatured: "Recomendado",
      planEquivalent: "Equivalente direto",
      planBtnFeat: "Assinar Pro",
      planBtnNormal: "Escolher",
      proofTitle: "Devs brasileiros<br />respirando IA o dia inteiro.",
      proofSubtitle: "Quem usa NexusMind está desenvolvendo em ritmo acelerado, criando agentes e automações. Não dá para parar.",
      faqTitle: "Dúvidas que<br />aparecem no DM.",
      faqSubtitle: "Se a sua não está aqui, nos chame no WhatsApp para atendimento direto. Sem bots.",
      finalTitle: "Pronto para rotear seus tokens?",
      finalSubtitle: "O Pix confirma em segundos. Em 4 minutos você está rodando o Claude Opus 4.7 no seu Cursor, pagando em real.",
      ctaWhatsApp: "Falar no WhatsApp",
      footerDesc: "Gateway de roteamento e cota para APIs de IA. Acesse Claude, GPT e Gemini com uma chave só, em real, no Pix. Não somos a Anthropic — somos seu intermediário brasileiro.",
      footerProduct: "Produto",
      footerChannels: "Canais",
      legalTerms: "Termos de uso",
      legalPrivacy: "Privacidade",
      legalRefunds: "Reembolso",
      footerBottomText: "Pagamento facilitado via Pix ou cartão",
      floatBtn: "Assinar · R$ 249",
      freeTrialTitle: "Experimente grátis hoje mesmo.",
      freeTrialDesc: "Crie sua conta em 30 segundos e ganhe US$ 5,00 em créditos de cortesia para testar os modelos de fronteira (Opus, o1, R1). Sem burocracia, sem cartão de crédito.",
      freeTrialBtn: "Experimente Grátis"
    }
  } else if (isZh.value) {
    return {
      navHow: "工作原理",
      navMath: "定价数学",
      navPricing: "价格",
      navFaq: "常见问题",
      btnDashboard: "控制台",
      btnLogin: "登录",
      btnStart: "开始使用",
      heroTitle: "一把钥匙。<br /><span class='nm-accent'>所有 AI 模型</span><br />即刻启用。",
      heroLede: "针对 Claude、GPT 和 Gemini 的智能路由 —— 实时计费，支持 Pix/信用卡。用于监控 Token、成本并在供应商发生故障时自动切换的控制面板。",
      ctaGetStarted: "获取我的 Key",
      ctaOpenDashboard: "打开控制台",
      ctaHowItWorks: "工作原理",
      heroPill1: "3 个供应商 · 14 个模型",
      heroPill2: "99.94% 在线率 · 30天",
      heroPill3: "★ 4.9 · 100+ 真实评价",
      consoleLive: "智能路由中",
      consoleRouting: "当前路由",
      metaReqs: "请求/秒",
      metaLatency: "p95 延迟",
      metaSavings: "平均节省",
      metaMonth: "月",
      tutTitle: "90 秒视频演示。",
      tutSubtitle: "从 Pix 支付到首次 API 调用。无多余旁白 —— 纯操作录屏演示。",
      tutNoAudio: "无音频",
      tutHint: "在代码中的 <code>videoEmbedUrl</code> 字段填入视频链接。",
      howTitle: "只需 3 步。<br />简单快捷。",
      howSubtitle: "您无需进行复杂的平台迁移，无需注册国外供应商账号，也不用付美元。只需更换 API 基址 (Base URL)。",
      step1Num: "创建 KEY",
      step1Title: "付款并获取 Key",
      step1Desc: "Pix 实时到账确认。付款后 Key 立即显示在您的控制面板，并发送至您的购买邮箱。",
      step2Num: "配置环境",
      step2Title: "更换 API 基址",
      step2Desc: "将 Key 和基址粘贴到您的 Cursor、Claude Code、脚本或 SDK 中。兼容任何 OpenAI 客户端样式。",
      step3Num: "开始使用",
      step3Title: "发送请求即刻路由",
      step3Desc: "请求到达路由网关后，我们将自动为您挑选当前延迟最低且成本最低的可用供应商响应。",
      mathTitle: "为什么只需 <span class='nm-accent'>R$ 249</span><br />而不是 R$ 1.020?",
      mathSubtitle: "没有任何隐藏条款。相同的模型，相同的 API，相同的性能。差额来自三个核心因素。",
      mathCol1: "月度订阅",
      mathCol2: "官方直签",
      mathCol3: "通过",
      mathColTotal: "每月总计",
      mathIncluded: "已包含",
      intTitle: "在您现有的<br />开发工具中接入。",
      intSubtitle: "修改 API 基址和 Key 即可无缝运行。支持任何兼容 OpenAI API 的客户端 —— 无需重构代码。",
      intTested: "已测试",
      intMoreTitle: "同时支持运行于",
      intMoreText: "Zed · Vercel AI SDK · LlamaIndex · OpenRouter 客户端 · Goose · Roo Code · Cody · Python/Node 脚本 — 任何兼容 OpenAI 密钥格式的工具。",
      provTitle: "同一密钥支持<br />超过 100 个模型。",
      provSubtitle: "首页无需列出所有模型：关键是 FragaAi 支持在高端、编程、推理、多模态和开源系列模型之间进行智能路由，且兼容 OpenAI/Anthropic 接口。",
      provFooterText: "<b>动态目录：</b>除了上述模型外，网关还提供 balanced-load, high-availability, coding, reasoning, economy, opus, flash 和 raiko 等组合。新模型进入路由器，无需您更改集成。",
      priceTitle: "三种方案。<br />无隐藏套路。",
      priceSubtitle: "随时取消。支持 Pix 或信用卡便捷支付。",
      planFeatured: "推荐方案",
      planEquivalent: "官方直签对应价格",
      planBtnFeat: "订阅 Pro 版",
      planBtnNormal: "选择",
      proofTitle: "众多开发者<br />全天候的高频使用。",
      proofSubtitle: "NexusMind 的用户遍布重构开发、自主代理和循环自动化流程。高稳定性保障工作不中断。",
      faqTitle: "常见私信问题。",
      faqSubtitle: "如果没有解答您的问题，欢迎点击 WhatsApp 联系人工客服。无机器人套路。",
      finalTitle: "准备好优化您的 Token 消耗了吗？",
      finalSubtitle: "付款在数秒内确认。只需 4 分钟即可在您的 Cursor 中运行 Claude Opus 4.7，支付本地货币。",
      ctaWhatsApp: "联系 WhatsApp",
      footerDesc: "智能路由与 Token 额度网关。一键使用 Claude、GPT 和 Gemini。我们不是 Anthropic 官方 —— 我们是您可信赖的本地代理商。",
      footerProduct: "产品",
      footerChannels: "渠道",
      legalTerms: "服务条款",
      legalPrivacy: "隐私政策",
      legalRefunds: "退款政策",
      footerBottomText: "支持 Pix 或信用卡付款",
      floatBtn: "订阅 · R$ 249",
      freeTrialTitle: "立即免费开始体验。",
      freeTrialDesc: "在 30 秒内创建您的账户，即可获得 5.00 美元的赠送额度，用于测试前沿模型（Opus, o1, R1）。无风险，无需绑定信用卡。",
      freeTrialBtn: "免费体验"
    }
  } else {
    return {
      navHow: "how_it_works",
      navMath: "math",
      navPricing: "price",
      navFaq: "faq",
      btnDashboard: "Dashboard",
      btnLogin: "Login",
      btnStart: "Get Started",
      heroTitle: "One key.<br /><span class='nm-accent'>All AI models</span><br />at your disposal.",
      heroLede: "Smart routing for Claude, GPT, and Gemini — in local currency, with Pix/Card. Dashboard to monitor tokens, cost, and automatic failover when a provider stalls.",
      ctaGetStarted: "Get my key now",
      ctaOpenDashboard: "Open Dashboard",
      ctaHowItWorks: "How it works",
      heroPill1: "3 providers · 14 models",
      heroPill2: "99.94% uptime · 30d",
      heroPill3: "★ 4.9 · 100+ reviews",
      consoleLive: "ROUTINGACTIVE",
      consoleRouting: "now routing",
      metaReqs: "req/s",
      metaLatency: "p95 latency",
      metaSavings: "avg savings",
      metaMonth: "mo",
      tutTitle: "Watch in 90 seconds.",
      tutSubtitle: "From Pix deposit to the first API call. No fluff, no audio — just the raw clicks on the screen.",
      tutNoAudio: "NO AUDIO",
      tutHint: "Paste your video link inside the <code>videoEmbedUrl</code> variable in code.",
      howTitle: "In 3 steps.<br />No nonsense.",
      howSubtitle: "No migration, no foreign accounts, no international transaction fees. Just change the base URL.",
      step1Num: "CREATE KEY",
      step1Title: "Pay and receive your key",
      step1Desc: "Instantly confirmed. The key appears in your dashboard and is sent to your purchase email.",
      step2Num: "CONFIGURE",
      step2Title: "Swap the base URL",
      step2Desc: "Paste the key and the Base URL in Cursor, Claude Code, custom script or SDK. Works with any OpenAI-compatible client.",
      step3Num: "USE",
      step3Title: "Send request, we route it",
      step3Desc: "Your request hits the router, we pick the cheapest healthy provider for your model, and serve the response.",
      mathTitle: "Why pay <span class='nm-accent'>R$ 249</span><br />instead of R$ 1,020?",
      mathSubtitle: "No catches. Same models, same API, same performance. The difference comes down to three factors.",
      mathCol1: "Monthly subscription",
      mathCol2: "Directly at provider",
      mathCol3: "Via",
      mathColTotal: "Total / month",
      mathIncluded: "included",
      intTitle: "Paste your key where<br />you already work.",
      intSubtitle: "Just swap the base URL and key, and it works. Any OpenAI-compatible client is supported — no code refactoring.",
      intTested: "Tested",
      intMoreTitle: "also runs on",
      intMoreText: "Zed · Vercel AI SDK · LlamaIndex · OpenRouter clients · Goose · Roo Code · Cody · Python/Node scripts — any tool accepting an OpenAI-compatible key format.",
      provTitle: "Over 100 models<br />on the same key.",
      provSubtitle: "The homepage doesn't need to list everything: the point is that FragaAi routes between premium, coding, reasoning, multimodal and open-source families via a compatible OpenAI/Anthropic base.",
      provFooterText: "<b>Dynamic catalog:</b> in addition to the models above, the gateway exposes combos such as balanced-load, high-availability, coding, reasoning, economy, opus, flash and raiko. New models enter the router without you changing integration.",
      priceTitle: "Three plans.<br />Zero fine print.",
      priceSubtitle: "Cancel anytime. Flexible payment via Pix or local cards.",
      planFeatured: "Recommended",
      planEquivalent: "Direct equivalent",
      planBtnFeat: "Get Pro Plan",
      planBtnNormal: "Choose",
      proofTitle: "Local developers<br />breathing AI all day long.",
      proofSubtitle: "Users are doing heavy refactoring, deploying autonomous agents, and running loops. Reliable uptime keeps them coding.",
      faqTitle: "Questions that<br />land in our DM.",
      faqSubtitle: "If yours isn't listed, ping us on WhatsApp for a direct talk. No automated bots.",
      finalTitle: "Ready to route your tokens?",
      finalSubtitle: "Confirmed in seconds. You will be running Claude Opus 4.7 in Cursor in less than 4 minutes.",
      ctaWhatsApp: "Talk on WhatsApp",
      footerDesc: "Quota and routing gateway for AI APIs. Access Claude, GPT, and Gemini with a single key. We are not Anthropic — we are your reliable local reseller.",
      footerProduct: "Product",
      footerChannels: "Channels",
      legalTerms: "Terms",
      legalPrivacy: "Privacy",
      legalRefunds: "Refund policy",
      footerBottomText: "Flexible payments via Pix or Card",
      floatBtn: "Subscribe · R$ 249",
      freeTrialTitle: "Start testing for free today.",
      freeTrialDesc: "Create your account in 30 seconds and get $5.00 in complimentary quickstart credits to test frontier models (Opus, o1, R1). No risk, no credit card required.",
      freeTrialBtn: "Try for Free"
    }
  }
})

// ── Routing Console Animation ──────────────────────────────
const activeRoute = ref(0)
const packetTick = ref(0)
interface LogEntry { id: number; t: string; path: string; prov: string; model: string; latency: number }
const logs = ref<LogEntry[]>([])
let routeInterval: ReturnType<typeof setInterval> | null = null

const nmRequests = [
  { y: 80, label: 'POST /v1/messages', method: 'POST /v1/messages', path: '/v1/messages', model: 'claude-opus-4.7', prov: 'AN_03' },
  { y: 140, label: 'POST /v1/chat/completions', method: 'POST /v1/chat/completions', path: '/v1/chat/completions', model: 'gpt-5', prov: 'OA_01' },
  { y: 200, label: 'POST /v1/embeddings', method: 'POST /v1/embeddings', path: '/v1/embeddings', model: 'gemini-2.5-flash', prov: 'GG_02' },
]
const nmProviders = [
  { y: 70, name: 'ANTHROPIC' },
  { y: 160, name: 'OPENAI' },
  { y: 250, name: 'GOOGLE' },
]

function tickRoute() {
  const next = (activeRoute.value + 1) % 3
  activeRoute.value = next
  packetTick.value++
  const now = new Date()
  const t = `${String(now.getHours()).padStart(2,'0')}:${String(now.getMinutes()).padStart(2,'0')}:${String(now.getSeconds()).padStart(2,'0')}`
  const latency = 280 + Math.floor(Math.random() * 220)
  const req = nmRequests[next]
  const entry: LogEntry = { id: now.getTime() + Math.random(), t, path: req.path, prov: req.prov, model: req.model, latency }
  logs.value = [entry, ...logs.value].slice(0, 3)
}

// ── Floating CTA Scroll Trigger ────────────────────────────
const floatVisible = ref(false)
const floatDismissed = ref(false)
function onScroll() {
  if (floatDismissed.value) return
  floatVisible.value = window.scrollY > 700
}

// ── Video URL and Chapters ─────────────────────────────────
const videoEmbedUrl = ref('')

const translatedChapters = computed(() => {
  if (isPt.value) {
    return [
      { t: '00:00', txt: 'Pagar via Pix em 3 cliques' },
      { t: '00:15', txt: 'Receber a chave no painel' },
      { t: '00:30', txt: 'Configurar no Cursor' },
      { t: '00:55', txt: 'Configurar no Claude Code' },
      { t: '01:15', txt: 'Primeira chamada · Opus 4.7' },
      { t: '01:25', txt: 'Monitorar uso e custo' }
    ]
  } else if (isZh.value) {
    return [
      { t: '00:00', txt: '点击 3 下通过 Pix 付款' },
      { t: '00:15', txt: '在控制面板接收 Key' },
      { t: '00:30', txt: '在 Cursor 中配置基址' },
      { t: '00:55', txt: '在 Claude Code 中配置环境' },
      { t: '01:15', txt: '测试首次调用 · Opus 4.7' },
      { t: '01:25', txt: '监控 Token 使用与成本' }
    ]
  } else {
    return [
      { t: '00:00', txt: 'Pay via Pix in 3 clicks' },
      { t: '00:15', txt: 'Get key in your dashboard' },
      { t: '00:30', txt: 'Configure inside Cursor' },
      { t: '00:55', txt: 'Configure inside Claude Code' },
      { t: '01:15', txt: 'First call using Opus 4.7' },
      { t: '01:25', txt: 'Monitor API use and token cost' }
    ]
  }
})

// ── Integrations ───────────────────────────────────────────
const translatedIntegrations = computed(() => {
  if (isPt.value) {
    return [
      { name: 'Cursor', icon: '▸', where: 'Configurações → Modelos → Adicionar base URL', env: 'OPENAI_BASE_URL', tested: true },
      { name: 'Claude Code', icon: '⌘', where: 'Variável de ambiente do shell', env: 'ANTHROPIC_BASE_URL', tested: true },
      { name: 'Windsurf', icon: '≋', where: 'Configurações → AI → Custom endpoint', env: 'BASE_URL', tested: true },
      { name: 'Cline', icon: '/', where: 'Extensão VS Code → Configurações', env: 'OPENAI_BASE_URL', tested: true },
      { name: 'Continue.dev', icon: '→', where: '~/.continue/config.json', env: 'apiBase', tested: true },
      { name: 'Aider', icon: ':', where: 'Flag --openai-api-base', env: 'OPENAI_API_BASE', tested: true },
      { name: 'n8n', icon: '⬡', where: 'Nó AI / OpenAI → URL personalizada', env: 'baseURL', tested: true },
      { name: 'LangChain', icon: 'λ', where: 'ChatOpenAI(base_url=…)', env: 'base_url', tested: true },
    ]
  } else if (isZh.value) {
    return [
      { name: 'Cursor', icon: '▸', where: '设置 → 模型 → 添加基址 (Base URL)', env: 'OPENAI_BASE_URL', tested: true },
      { name: 'Claude Code', icon: '⌘', where: 'Shell 环境变量', env: 'ANTHROPIC_BASE_URL', tested: true },
      { name: 'Windsurf', icon: '≋', where: '设置 → AI → 自定义端点', env: 'BASE_URL', tested: true },
      { name: 'Cline', icon: '/', where: 'VS Code 插件 → 设置', env: 'OPENAI_BASE_URL', tested: true },
      { name: 'Continue.dev', icon: '→', where: '~/.continue/config.json', env: 'apiBase', tested: true },
      { name: 'Aider', icon: ':', where: '--openai-api-base 参数标志', env: 'OPENAI_API_BASE', tested: true },
      { name: 'n8n', icon: '⬡', where: 'AI / OpenAI 节点 → 自定义 URL', env: 'baseURL', tested: true },
      { name: 'LangChain', icon: 'λ', where: 'ChatOpenAI(base_url=…)', env: 'base_url', tested: true },
    ]
  } else {
    return [
      { name: 'Cursor', icon: '▸', where: 'Settings → Models → Add base URL', env: 'OPENAI_BASE_URL', tested: true },
      { name: 'Claude Code', icon: '⌘', where: 'Shell Env Variable', env: 'ANTHROPIC_BASE_URL', tested: true },
      { name: 'Windsurf', icon: '≋', where: 'Settings → AI → Custom endpoint', env: 'BASE_URL', tested: true },
      { name: 'Cline', icon: '/', where: 'VS Code Ext → Settings', env: 'OPENAI_BASE_URL', tested: true },
      { name: 'Continue.dev', icon: '→', where: '~/.continue/config.json', env: 'apiBase', tested: true },
      { name: 'Aider', icon: ':', where: '--openai-api-base flag', env: 'OPENAI_API_BASE', tested: true },
      { name: 'n8n', icon: '⬡', where: 'AI / OpenAI node → custom URL', env: 'baseURL', tested: true },
      { name: 'LangChain', icon: 'λ', where: 'ChatOpenAI(base_url=…)', env: 'base_url', tested: true },
    ]
  }
})

// ── Matemática ─────────────────────────────────────────────
const activeMath = ref(0)
const rcMathRows = [
  { label: 'Claude Pro / Max', direct: 'R$ 580' },
  { label: 'ChatGPT Plus / Pro', direct: 'R$ 320' },
  { label: 'Gemini Advanced', direct: 'R$ 120' },
]

const translatedMathPoints = computed(() => {
  if (isPt.value) {
    return [
      { tag: '01 · Fator', title: 'Compra em volume', body: 'Compramos pacotes de tokens dos provedores em larga escala. Volume continuado se traduz em desconto de atacado — e nós repassamos essa economia para você.' },
      { tag: '02 · Fator', title: 'Roteamento inteligente', body: 'Para cada chamada, selecionamos automaticamente a conta/provedor com menor custo e latência ideal. Você obtém sempre a melhor taxa de conversão.' },
      { tag: '03 · Fator', title: 'Sem múltiplas assinaturas', body: 'Em vez de assinar Claude Max + ChatGPT Pro + Gemini Ultra (R$ 1.020+), você paga uma única taxa mensal e consome todos conforme sua necessidade.' }
    ]
  } else if (isZh.value) {
    return [
      { tag: '01 · 因素', title: '大容量采购优惠', body: '我们从各供应商处大批量采购 Token。持续的采购量转化为可观的批发折扣 —— 并直接让利给我们的开发者。' },
      { tag: '02 · 因素', title: '智能负载路由', body: '每次请求发出时，系统自动挑选当前健康度最佳且单价最合算的账户路由。始终确保最高的性能开销比。' },
      { tag: '03 · 因素', title: '省去多平台订阅', body: '无需分别开通 Claude Max + ChatGPT Pro + Gemini Ultra 订阅（每月 R$ 1,020+），统一通过一个账户共享额度，按需灵活分配。' }
    ]
  } else {
    return [
      { tag: '01 · Factor', title: 'Bulk volume scaling', body: 'We secure massive token quotas directly from major providers. High volume unlocks lower rates — and we pass the savings directly to you.' },
      { tag: '02 · Factor', title: 'Intelligent load routing', body: 'For every API request, we instantly pick the cheapest healthy route. You run your models at optimized cost, completely automatically.' },
      { tag: '03 · Factor', title: 'No multiple subscription plans', body: 'Instead of keeping separate plans for Claude Max + ChatGPT Pro + Gemini Ultra, a single unified subscription covers all models in local currency.' }
    ]
  }
})

// ── Providers ─────────────────────────────────────────────
const translatedProviderCards = computed(() => {
  if (isPt.value) {
    return [
      {
        code: 'CL',
        name: 'Claude / Anthropic',
        family: 'Família Claude 4.x',
        modelList: ['Opus 4.7', 'Opus 4.6', 'Sonnet 4.6', 'Sonnet 4.5', 'Haiku 4.5', 'Thinking'],
        color: '#d97757',
      },
      {
        code: 'CX',
        name: 'GPT / Codex',
        family: 'Família GPT 5.x + coding',
        modelList: ['GPT-5.5', 'GPT-5.5 xHigh', 'GPT-5.5 High', 'GPT-5.4', 'GPT-5.4 Mini', 'Codex Spark', 'Auto Review'],
        color: '#10a37f',
      },
      {
        code: 'GM',
        name: 'Gemini / Google',
        family: 'Pro, Flash, Lite e multimodal',
        modelList: ['Gemini 3.1 Pro', 'Gemini 3 Pro', 'Gemini 3 Flash', 'Gemini 2.5 Pro', '2.5 Flash', 'Flash Lite', 'Image'],
        color: '#4285f4',
      },
      {
        code: 'KR',
        name: 'Kiro / Agentic',
        family: 'Modelos para agentes e IDE',
        modelList: ['Auto Kiro', 'Claude Opus 4.7', 'Claude Sonnet 4.6', 'DeepSeek 3.2', 'MiniMax M2.5', 'GLM-5', 'Qwen3 Coder'],
        color: '#8b5cf6',
      },
      {
        code: 'OS',
        name: 'Open-source / Reasoning',
        family: 'DeepSeek, Qwen, Llama, Kimi, GLM',
        modelList: ['DeepSeek V3.2', 'DeepSeek R1', 'Qwen3 235B', 'Qwen3 Coder', 'Llama 4 Scout', 'Kimi K2.5', 'GLM-4.7'],
        color: '#f97316',
      },
      {
        code: 'MX',
        name: 'Mistral / Groq / Extras',
        family: 'Velocidade, áudio e modelos europeus',
        modelList: ['Mistral Large', 'Mistral Medium', 'Codestral', 'Devstral', 'Groq Llama', 'GPT-OSS 120B', 'Whisper'],
        color: '#ec4899',
      },
    ]
  } else if (isZh.value) {
    return [
      {
        code: 'CL',
        name: 'Claude / Anthropic',
        family: 'Claude 4.x 系列',
        modelList: ['Opus 4.7', 'Opus 4.6', 'Sonnet 4.6', 'Sonnet 4.5', 'Haiku 4.5', 'Thinking'],
        color: '#d97757',
      },
      {
        code: 'CX',
        name: 'GPT / Codex',
        family: 'GPT 5.x 系列 + 编程专用',
        modelList: ['GPT-5.5', 'GPT-5.5 xHigh', 'GPT-5.5 High', 'GPT-5.4', 'GPT-5.4 Mini', 'Codex Spark', 'Auto Review'],
        color: '#10a37f',
      },
      {
        code: 'GM',
        name: 'Gemini / Google',
        family: 'Pro, Flash, Lite 与多模态',
        modelList: ['Gemini 3.1 Pro', 'Gemini 3 Pro', 'Gemini 3 Flash', 'Gemini 2.5 Pro', '2.5 Flash', 'Flash Lite', 'Image'],
        color: '#4285f4',
      },
      {
        code: 'KR',
        name: 'Kiro / Agentic',
        family: '智能体与 IDE 专用模型',
        modelList: ['Auto Kiro', 'Claude Opus 4.7', 'Claude Sonnet 4.6', 'DeepSeek 3.2', 'MiniMax M2.5', 'GLM-5', 'Qwen3 Coder'],
        color: '#8b5cf6',
      },
      {
        code: 'OS',
        name: 'Open-source / Reasoning',
        family: 'DeepSeek, Qwen, Llama, Kimi, GLM',
        modelList: ['DeepSeek V3.2', 'DeepSeek R1', 'Qwen3 235B', 'Qwen3 Coder', 'Llama 4 Scout', 'Kimi K2.5', 'GLM-4.7'],
        color: '#f97316',
      },
      {
        code: 'MX',
        name: 'Mistral / Groq / Extras',
        family: '高并发、语音及欧洲系列模型',
        modelList: ['Mistral Large', 'Mistral Medium', 'Codestral', 'Devstral', 'Groq Llama', 'GPT-OSS 120B', 'Whisper'],
        color: '#ec4899',
      },
    ]
  } else {
    return [
      {
        code: 'CL',
        name: 'Claude / Anthropic',
        family: 'Claude 4.x Family',
        modelList: ['Opus 4.7', 'Opus 4.6', 'Sonnet 4.6', 'Sonnet 4.5', 'Haiku 4.5', 'Thinking'],
        color: '#d97757',
      },
      {
        code: 'CX',
        name: 'GPT / Codex',
        family: 'GPT 5.x Family + coding',
        modelList: ['GPT-5.5', 'GPT-5.5 xHigh', 'GPT-5.5 High', 'GPT-5.4', 'GPT-5.4 Mini', 'Codex Spark', 'Auto Review'],
        color: '#10a37f',
      },
      {
        code: 'GM',
        name: 'Gemini / Google',
        family: 'Pro, Flash, Lite and multimodal',
        modelList: ['Gemini 3.1 Pro', 'Gemini 3 Pro', 'Gemini 3 Flash', 'Gemini 2.5 Pro', '2.5 Flash', 'Flash Lite', 'Image'],
        color: '#4285f4',
      },
      {
        code: 'KR',
        name: 'Kiro / Agentic',
        family: 'Agent & IDE specific models',
        modelList: ['Auto Kiro', 'Claude Opus 4.7', 'Claude Sonnet 4.6', 'DeepSeek 3.2', 'MiniMax M2.5', 'GLM-5', 'Qwen3 Coder'],
        color: '#8b5cf6',
      },
      {
        code: 'OS',
        name: 'Open-source / Reasoning',
        family: 'DeepSeek, Qwen, Llama, Kimi, GLM',
        modelList: ['DeepSeek V3.2', 'DeepSeek R1', 'Qwen3 235B', 'Qwen3 Coder', 'Llama 4 Scout', 'Kimi K2.5', 'GLM-4.7'],
        color: '#f97316',
      },
      {
        code: 'MX',
        name: 'Mistral / Groq / Extras',
        family: 'High-speed, audio & European models',
        modelList: ['Mistral Large', 'Mistral Medium', 'Codestral', 'Devstral', 'Groq Llama', 'GPT-OSS 120B', 'Whisper'],
        color: '#ec4899',
      },
    ]
  }
})

// ── Plans ─────────────────────────────────────────────────
const translatedPlans = computed(() => {
  if (isPt.value) {
    return [
      {
        type: 'Starter', name: 'Solo dev', feat: false, price: 79, oldPrice: null,
        desc: 'Para quem está iniciando ou usando inteligência artificial em projetos de escopo pessoal.',
        features: ['1 chave de API ativa', 'Acesso a Claude, GPT e Gemini', 'Limite mensal de 50M tokens', 'Painel de uso em tempo real', 'Suporte padrão por e-mail'],
        muted: ['Roteamento prioritário', 'Gestão multi-times'],
      },
      {
        type: 'Recomendado', name: 'Pro', feat: true, price: 249, oldPrice: 1020,
        desc: 'Acesso total a todos os modelos premium de fronteira, sem rodeios ou limites de cota restritivos.',
        features: ['Chaves de API ilimitadas por projeto', 'Acesso completo a todos os modelos premium', 'Limite mensal de 500M tokens', 'Roteamento prioritário entre clusters', 'Failover de latência automática integrado', 'Logs completos e exportáveis', 'Suporte direto via WhatsApp corporativo'],
        muted: [],
      },
      {
        type: 'Team', name: 'Equipe', feat: false, price: 749, oldPrice: null,
        desc: 'Para times de engenharia de software e operações baseadas em agentes em alta escala.',
        features: ['Todos os recursos do plano Pro inclusos', 'Grupos e quotas segmentadas por time', 'Auditoria completa de token usage', 'Suporte de integração avançada (SSO)', 'SLA contratual de 99.9% de uptime', 'Onboarding dedicado com engenheiro B2B'],
        muted: [],
      }
    ]
  } else if (isZh.value) {
    return [
      {
        type: 'Starter', name: '个人开发者', feat: false, price: 79, oldPrice: null,
        desc: '适用于学习试用或轻度个人项目的开发者环境。',
        features: ['1 个激活的 API 密钥', '支持接入 Claude、GPT 和 Gemini', '每月 50M Token 额度', '实时使用额度统计图表', '标准的电子邮件技术支持'],
        muted: ['优先队列智能路由', '多团队联合管理'],
      },
      {
        type: '推荐方案', name: 'Pro 专业版', feat: true, price: 249, oldPrice: 1020,
        desc: '全功能无限制接入所有前沿模型，无 5 小时高频限流困扰。',
        features: ['按项目创建无限个 API 密钥', '所有前沿模型不限版本接入', '每月高达 500M Token 额度', '网关优先通道路由保障', '智能故障瞬时自动路由切换', '详尽的使用日志与报表导出', '专属 WhatsApp/微信 及时客户经理支持'],
        muted: [],
      },
      {
        type: 'Team', name: '团队企业版', feat: false, price: 749, oldPrice: null,
        desc: '适用于高频应用、软件工程团队或中大型自动化生产流程。',
        features: ['包含 Pro 专业版全部权益', '团队内部分级配额与预算隔离', 'Token 审计与细粒度使用报告', '企业级单点登录 (SSO) 与安全审计', '99.9% 稳定率 SLA 合同条款保障', '资深 B2B 架构师对接导入'],
        muted: [],
      }
    ]
  } else {
    return [
      {
        type: 'Starter', name: 'Solo dev', feat: false, price: 79, oldPrice: null,
        desc: 'For individuals starting out or running light workloads for personal AI sandboxes.',
        features: ['1 active API key', 'Access to Claude, GPT, and Gemini', '50M tokens monthly limit', 'Real-time usage dashboard', 'Standard email support'],
        muted: ['Priority smart routing queue', 'Multi-team hierarchy management'],
      },
      {
        type: 'Recommended', name: 'Pro', feat: true, price: 249, oldPrice: 1020,
        desc: 'Full frontier access to premium models with zero restrictive 5-hour rate limits.',
        features: ['Unlimited project-isolated API keys', 'Access to all available premium models', '500M tokens monthly limit', 'Priority cluster routing paths', 'Built-in instant failover mechanics', 'Full usage statistics export option', 'Direct Slack/WhatsApp business support'],
        muted: [],
      },
      {
        type: 'Team', name: 'Team', feat: false, price: 749, oldPrice: null,
        desc: 'For software engineering teams scaling up agentic workflows and integrations.',
        features: ['Includes all Pro plan benefits', 'Isolated budgets and limits per team member', 'Advanced Token usage audits', 'SSO authentication support', 'Contractual 99.9% SLA uptime guarantee', 'Dedicated onboarding session with B2B architect'],
        muted: [],
      }
    ]
  }
})

// ── Testimonials ──────────────────────────────────────────
const translatedTestimonials = computed(() => {
  if (isPt.value) {
    return [
      { quote: 'Migrei 3 produtos para a NexusMind em um final de semana. Economia de R$ 2.300/mês e parou de cair quando os provedores principais engasgam.', name: 'Marina V.', role: 'Tech Lead · Fintech SP', initial: 'MV', meta: ['ECONOMIA', 'R$ 2.300/MÊS'] },
      { quote: 'O setup levou menos de 4 minutos. Troquei a base URL no meu Cursor e na hora estava rodando com o Opus 4.7.', name: 'Caio P.', role: 'Indie Hacker · POA', initial: 'CP', meta: ['SETUP', '4 MINUTOS'] },
      { quote: 'O failover inteligente salvou meu deploy. A Anthropic retornou erro 529 quatro vezes esse mês e nenhum usuário percebeu.', name: 'Heitor M.', role: 'CTO · Agência SaaS', initial: 'HM', meta: ['ESTABILIDADE', '99.94%'] },
    ]
  } else if (isZh.value) {
    return [
      { quote: '我在周末将 3 个项目迁移到了 NexusMind。每月直接节省 R$ 2,300，且再也没遇到过因官方服务波动导致的崩溃。', name: 'Marina V.', role: '技术负责人 · 圣保罗金融科技', initial: 'MV', meta: ['降本增效', 'R$ 2,300/月'] },
      { quote: '配置流程只花了不到 4 分钟。在 Cursor 中填入新的 API 基址，立即成功调用了最新的 Opus 4.7 模型。', name: 'Caio P.', role: '独立开发者 · 阿雷格里港', initial: 'CP', meta: ['部署用时', '4 分钟'] },
      { quote: '智能路由故障转移拯救了我们的生产线。当月官方服务发生 4 次大面积抖动，我们的客户却完全没有感知。', name: 'Heitor M.', role: '首席技术官 · 软件服务商', initial: 'HM', meta: ['稳定运行', '99.94%'] },
    ]
  } else {
    return [
      { quote: 'Moved 3 projects to NexusMind in one afternoon. Saving R$ 2,300/mo and no more production outages when primary providers go down.', name: 'Marina V.', role: 'Tech Lead · Fintech SP', initial: 'MV', meta: ['SAVINGS', 'R$ 2,300/MO'] },
      { quote: 'Setup took me less than 4 minutes. Swapped the base URL in my Cursor settings and immediately had Opus 4.7 working.', name: 'Caio P.', role: 'Indie Hacker · POA', initial: 'CP', meta: ['SETUP', '4 MINS'] },
      { quote: 'Smart failover saved our deployment. Anthropic had four separate 529 outages this month, and none of our users noticed a thing.', name: 'Heitor M.', role: 'CTO · SaaS Agency', initial: 'HM', meta: ['UPTIME', '99.94%'] },
    ]
  }
})

// ── FAQ ───────────────────────────────────────────────────
const openFaq = ref(0)
function toggleFaq(i: number) {
  openFaq.value = openFaq.value === i ? -1 : i
}

const translatedFaqs = computed(() => {
  if (isPt.value) {
    return [
      { q: 'É a mesma API oficial dos provedores?', a: 'Sim. O NexusMind atua como um gateway de roteamento e cota. Suas requisições são encaminhadas de ponta a ponta para os endpoints oficiais da Anthropic, OpenAI e Google. Você utiliza os mesmos SDKs oficiais, apenas trocando a Base URL e a chave de autenticação.' },
      { q: 'Como é cobrado? Por token ou por assinatura mensal?', a: 'O faturamento é mensal com pacotes de tokens inclusos. Isso evita flutuações de custos devido a oscilações no câmbio do dólar ou picos inesperados de uso. Caso sua equipe consuma a cota inclusa no período, você será notificado com antecedência para optar pelo upgrade.' },
      { q: 'Em quais IDEs e ferramentas posso integrar a chave?', a: 'Em qualquer cliente ou biblioteca que aceite uma chave no padrão OpenAI ou Anthropic: Cursor, Claude Code, Cline, Aider, Continue, n8n, LangChain, além de scripts personalizados em Python, Node.js, Go e Rust. Basta trocar a URL de destino e colar sua chave.' },
      { q: 'A plataforma oferece garantia de reembolso?', a: 'Sim, oferecemos garantia de arrependimento de 7 dias com reembolso integral de 100%, sem burocracias. Válido para pagamentos via Pix ou cartão.' },
      { q: 'Posso compartilhar minha chave com múltiplos computadores?', a: 'Sim. A chave pode ser associada a múltiplos dispositivos do mesmo desenvolvedor ou equipe. Monitoramos padrões incomuns para evitar abusos de compartilhamento público de credenciais.' },
      { q: 'Vocês são parceiros oficiais da Anthropic/OpenAI?', a: 'Somos um intermediário e agregador independente sediado no Brasil. Adquirimos volumes corporativos e revendemos em conformidade com as políticas comerciais de volume dos respectivos provedores.' },
      { q: 'Quanto tempo leva para minha chave ser ativada após a compra?', a: 'A ativação é automática. Pagamentos via Pix são compensados instantaneamente. Cartões de crédito levam em média até 2 minutos para processamento e liberação no painel.' },
      { q: 'O que acontece se a API da Anthropic cair globalmente?', a: 'Nosso roteador (Nexus Engine) aciona um failover dinâmico instantâneo. Se a chamada para o modelo desejado falhar ou retornar códigos de erro de servidor (5xx), a requisição é redirecionada para outra conta saudável em milissegundos.' },
    ]
  } else if (isZh.value) {
    return [
      { q: 'NexusMind 的 API 是否和官方一致？', a: '完全一致。NexusMind 是一个智能网关层，所有请求均直接通过加密通道转发至 Anthropic、OpenAI 和 Google 的官方 API 接口。您仍然可以使用官方提供的 SDK，只需修改基准 URL (Base URL) 并填入我们的 API Key 即可。' },
      { q: '它是按 Token 计费还是按月订阅？', a: '采用包月计费模式，每个档位的套餐均包含了足额的 Token。这使得企业免受汇率波动或团队意外高频调用导致的突发巨额账单困扰。如果超出额度，系统会提前通知，您可以随时平滑升级。' },
      { q: '它支持哪些 IDE 和开发框架？', a: '支持市面上所有兼容 OpenAI 接口格式的软件及框架。例如 Cursor, Claude Code, Cline, Aider, Continue, n8n 或者是您使用 Python, Node.js 编写的自定义调用脚本。' },
      { q: '是否支持退款保障？', a: '支持。我们提供 7 天无理由全额退款保障。如果您对服务有任何不满意，联系客服即可快速办理 Pix 或信用卡原路退款。' },
      { q: '一个 Key 可以在多台电脑上同时使用吗？', a: '可以。开发密钥支持绑定在同一个开发者或团队的多台设备上。我们仅对恶意的公开共享 Key 进行防滥用监测，正常的日常多机联调不受任何限制。' },
      { q: '你们是 Anthropic/OpenAI 的官方代理商吗？', a: '我们是位于巴西的独立技术网关提供商。我们通过大宗企业协议采购流量包，并在合规框架下以巴西本地货币和更符合本地生态的方式提供分销支持。' },
      { q: '购买后需要多久才能开通使用？', a: '系统实时自动激活。Pix 付款在数秒内即可到账确认并激活 Key；信用卡支付通常在 2 分钟内完成安全审核并分配资源。' },
      { q: '如果 Anthropic 等官方服务器崩溃了该怎么办？', a: '我们的智能路由系统内置了主动故障转移机制。当主通道发生超时或返回 5xx 错误时，网关将在毫秒级将流量重新分发至其他健康的备用通道，保障您的工作流程不受影响。' },
    ]
  } else {
    return [
      { q: 'Is this the same official API from the providers?', a: 'Yes. NexusMind operates as an API routing and quota gateway. Your calls travel directly to the official Anthropic, OpenAI, and Google endpoints. You can use all official SDKs without modification, simply changing the target Base URL and authentication key.' },
      { q: 'How is it billed? By token consumption or by monthly subscription?', a: 'We bill monthly with generous token quotas included. This keeps your monthly costs predictable and eliminates exchange rate spikes or billing shocks. If you exceed the quota, you will be notified to scale up.' },
      { q: 'Which IDEs and tools can I integrate the key into?', a: 'Any client or library that is compatible with OpenAI or Anthropic API formats: Cursor, Claude Code, Cline, Aider, Continue, n8n, LangChain, as well as custom Python, Node.js, and Go scripts.' },
      { q: 'Does the platform offer a refund guarantee?', a: 'Yes, we provide a 7-day money-back guarantee with a 100% refund. Valid for all Pix and local credit card subscriptions.' },
      { q: 'Can I share my key across multiple devices?', a: 'Yes. The key can be used across multiple devices belonging to the same developer or team. We monitor credentials for public misuse to prevent security locks.' },
      { q: 'Are you officially associated with Anthropic or OpenAI?', a: 'We are an independent technology reseller based in Brazil. We acquire corporate volume commitments and package them for local businesses in compliance with standard volume term sheets.' },
      { q: 'How long does key activation take after purchase?', a: 'Activation is instant. Pix transactions confirm in seconds. Credit cards take up to 2 minutes on average for payment gateway clearing.' },
      { q: 'What happens if Anthropic has a global outage?', a: 'Our router (Nexus Engine) activates zero-latency failover. If a request to your model times out or returns a 5xx server error, the call is instantly redirected to a healthy cluster node.' },
    ]
  }
})

// Scroll reveal observer
let revealObserver: IntersectionObserver | null = null
function setupReveal() {
  if (typeof IntersectionObserver === 'undefined') return
  revealObserver = new IntersectionObserver((entries) => {
    entries.forEach(e => {
      if (e.isIntersecting) {
        e.target.classList.add('nm-reveal-in')
        revealObserver?.unobserve(e.target)
      }
    })
  }, { threshold: 0.12, rootMargin: '0px 0px -60px 0px' })
  document.querySelectorAll('.nm-reveal').forEach(el => revealObserver?.observe(el))
}

onMounted(() => {
  tickRoute()
  routeInterval = setInterval(tickRoute, 1600)
  window.addEventListener('scroll', onScroll, { passive: true })
  onScroll()
  setTimeout(setupReveal, 50)
})

onUnmounted(() => {
  if (routeInterval) clearInterval(routeInterval)
  window.removeEventListener('scroll', onScroll)
  if (revealObserver) revealObserver.disconnect()
})
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=Outfit:wght@500;600;700;800&family=Syne:wght@600;700;800&display=swap');

.nm-home {
  /* Obsidian Space Core Colors */
  --nm-bg:      #030712;
  --nm-bg-1:    #0b0f19;
  --nm-bg-2:    #111827;
  --nm-line:    rgba(255, 255, 255, 0.04);
  --nm-line-2:  rgba(255, 255, 255, 0.08);
  --nm-fg:      #f9fafb;
  --nm-fg-1:    #e5e7eb;
  --nm-fg-2:    #9ca3af;
  --nm-fg-3:    #4b5563;

  /* Dashboard Theme Colors: Teal / Cyan */
  --nm-teal:      #14b8a6;
  --nm-teal-dark: #0f766e;
  --nm-violet:    #7c3aed;
  --nm-emerald:   #10b981;
  --nm-gray:      #374151;

  --nm-mono:    "JetBrains Mono", ui-monospace, SFMono-Regular, Menlo, monospace;
  --nm-sans:    "Inter", system-ui, -apple-system, sans-serif;
  --nm-title:   "Outfit", sans-serif;
  --nm-r:       12px;

  background: var(--nm-bg);
  color: var(--nm-fg);
  font-family: var(--nm-sans);
  font-feature-settings: "ss01","cv11";
  -webkit-font-smoothing: antialiased;
  min-height: 100vh;
  overflow-x: hidden;
  position: relative;
}

/* ── Grid & Glow layers ── */
.nm-bg-grid {
  position: fixed; inset: 0; z-index: 0; pointer-events: none;
  background-image: radial-gradient(rgba(20, 184, 166, 0.03) 1px, transparent 1px);
  background-size: 32px 32px;
  mask-image: radial-gradient(ellipse at 50% 30%, black 20%, transparent 80%);
}
.nm-bg-glow {
  position: fixed; inset: 0; z-index: 0; pointer-events: none;
  background:
    radial-gradient(800px 400px at 80% -10%, rgba(20, 184, 166, 0.12), transparent 60%),
    radial-gradient(600px 400px at 0% 30%, rgba(124, 58, 237, 0.08), transparent 60%);
}

.nm-wrap { max-width: 1240px; margin: 0 auto; padding: 0 32px; position: relative; z-index: 1; }
@media (max-width: 720px) { .nm-wrap { padding: 0 20px; } }

/* ── Header ── */
.nm-header {
  position: sticky; top: 0; z-index: 50;
  backdrop-filter: blur(14px);
  background: rgba(3, 7, 18, 0.75);
  border-bottom: 1px solid var(--nm-line);
}
.nm-header-inner { display: flex; align-items: center; justify-content: space-between; height: 68px; }
.nm-brand { display: flex; align-items: center; gap: 12px; text-decoration: none; }
.nm-brand-logo {
  width: 32px; height: 32px; border-radius: 8px; display: grid; place-items: center;
  background: #090d16; border: 1px solid var(--nm-line-2); overflow: hidden;
}
.nm-brand-logo img { width: 28px; height: 28px; object-fit: contain; }
.nm-brand-name { display: flex; flex-direction: column; line-height: 1.05; }
.nm-brand-name b { font-weight: 700; font-family: var(--nm-title); font-size: 15px; letter-spacing: -0.01em; color: var(--nm-fg); }
.nm-brand-name span { font-family: var(--nm-mono); font-size: 10px; color: var(--nm-fg-3); letter-spacing: 0.04em; text-transform: uppercase; }

.nm-nav { display: flex; align-items: center; gap: 24px; }
@media (max-width: 860px) { .nm-nav { display: none; } }
.nm-nav-link { font-family: var(--nm-mono); font-size: 12px; color: var(--nm-fg-2); letter-spacing: 0.04em; text-transform: uppercase; transition: color .15s; text-decoration: none; }
.nm-nav-link:hover { color: var(--nm-teal); }
.nm-actions { display: flex; align-items: center; gap: 12px; }

/* ── Buttons ── */
.nm-btn {
  display: inline-flex; align-items: center; gap: 8px;
  height: 40px; padding: 0 18px; border-radius: 8px;
  font-family: var(--nm-mono); font-size: 12px; font-weight: 600;
  letter-spacing: 0.04em; text-transform: uppercase;
  border: 1px solid transparent; cursor: pointer; transition: all .15s;
  white-space: nowrap; text-decoration: none;
}
.nm-btn-lg { height: 48px; padding: 0 24px; font-size: 13px; }
.nm-btn-primary {
  background: var(--nm-teal); color: #020617; border-color: var(--nm-teal);
  box-shadow: 0 0 15px rgba(20, 184, 166, 0.2);
}
.nm-btn-primary:hover {
  background: #2dd4bf; border-color: #2dd4bf;
  box-shadow: 0 0 25px rgba(20, 184, 166, 0.35);
  transform: translateY(-1px);
}
.nm-btn-ghost { background: transparent; color: var(--nm-fg); border-color: var(--nm-line-2); }
.nm-btn-ghost:hover { border-color: var(--nm-teal); color: var(--nm-teal); }
.nm-btn-dark { background: #0c121e; color: var(--nm-fg); border-color: var(--nm-line-2); }
.nm-btn-dark:hover { background: #111827; border-color: var(--nm-line-2); }
.nm-arr { font-family: var(--nm-sans); }

/* ── Sections ── */
.nm-section { padding: 96px 0; }
@media (max-width: 720px) { .nm-section { padding: 64px 0; } }
.border-t { border-top: 1px solid var(--nm-line); }

.nm-s-label {
  display: inline-flex; align-items: center; gap: 8px;
  font-family: var(--nm-mono); font-size: 11px; letter-spacing: 0.14em;
  color: var(--nm-teal); text-transform: uppercase; margin-bottom: 20px;
}
.nm-sq { width: 6px; height: 6px; background: var(--nm-teal); display: inline-block; }
.nm-tic { color: var(--nm-fg-3); }
.nm-s-title {
  font-family: var(--nm-title); font-size: clamp(32px, 4.2vw, 54px); font-weight: 800;
  line-height: 1.06; letter-spacing: -0.02em; margin: 0 0 18px;
  text-wrap: balance;
}
.nm-s-sub { font-size: 16px; color: var(--nm-fg-2); max-width: 640px; line-height: 1.6; margin: 0; }
.nm-accent { color: var(--nm-teal); background: linear-gradient(to right, var(--nm-teal), var(--nm-violet)); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
.nm-strike { text-decoration: line-through; color: var(--nm-fg-3); }

/* ── Hero ── */
.nm-hero { padding-top: 80px; padding-bottom: 96px; }
.nm-hero-grid { display: grid; grid-template-columns: 1.15fr 1fr; gap: 48px; align-items: center; }
@media (max-width: 1040px) { .nm-hero-grid { grid-template-columns: 1fr; gap: 56px; } }

.nm-eyebrow {
  display: inline-flex; align-items: center; gap: 8px;
  font-family: var(--nm-mono); font-size: 11px; text-transform: uppercase; letter-spacing: 0.12em;
  color: var(--nm-fg-1); background: #0b0f19; border: 1px solid var(--nm-line-2);
  padding: 6px 14px; border-radius: 99px; margin-bottom: 20px;
}
.nm-eyebrow-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--nm-teal); box-shadow: 0 0 8px var(--nm-teal); }
.nm-eyebrow-v { color: var(--nm-violet); font-weight: 600; }
.nm-h1 { font-size: clamp(42px,5.8vw,72px); font-weight: 800; line-height: 1.02; letter-spacing: -0.03em; margin: 0 0 22px; text-wrap: balance; font-family: var(--nm-title); }
.nm-lede { font-size: 18px; line-height: 1.55; color: var(--nm-fg-1); max-width: 540px; margin: 0 0 32px; }
.nm-cta-row { display: flex; flex-wrap: wrap; gap: 14px; margin-bottom: 32px; }
.nm-pill-row { display: flex; flex-wrap: wrap; gap: 20px; font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-2); }
.nm-pill-item { display: inline-flex; align-items: center; gap: 8px; }
.nm-dot-green { width: 8px; height: 8px; border-radius: 50%; background: var(--nm-emerald); box-shadow: 0 0 8px var(--nm-emerald); display: inline-block; }
.nm-ok { color: var(--nm-emerald); }

/* ── Console UI ── */
.nm-console {
  background: linear-gradient(180deg, #090c15 0%, #05080f 100%);
  border: 1px solid var(--nm-line-2); border-radius: var(--nm-r); overflow: hidden;
  box-shadow: 0 15px 40px rgba(0,0,0,0.5), inset 0 1px 0 rgba(255,255,255,0.05);
}
.nm-console-head {
  display: flex; align-items: center; gap: 10px; height: 38px; padding: 0 16px;
  border-bottom: 1px solid var(--nm-line); background: #05080f;
  font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-3);
}
.nm-console-dots { display: inline-flex; gap: 6px; margin-right: 8px; }
.nm-console-dots i { width: 10px; height: 10px; border-radius: 50%; display: inline-block; }
.nm-console-dots i.red { background: #ef4444; }
.nm-console-dots i.yellow { background: #f59e0b; }
.nm-console-dots i.green { background: #10b981; }
.nm-console-title { color: var(--nm-fg-2); }
.nm-console-live { margin-left: auto; display: inline-flex; align-items: center; gap: 6px; color: var(--nm-emerald); font-weight: 600; letter-spacing: 0.05em; }
.nm-live-dot { width: 7px; height: 7px; border-radius: 50%; background: var(--nm-emerald); box-shadow: 0 0 10px var(--nm-emerald); display: inline-block; animation: pulse 1.5s infinite; }
.nm-console-body { padding: 20px; }

/* SVG Lanes & Network */
.nm-diagram { width: 100%; height: 300px; display: block; }
.nm-lane { fill: none; stroke: url(#nmG1); stroke-width: 1.5; opacity: .15; stroke-dasharray: 4 6; }
.nm-lane-active { opacity: 1; stroke-dasharray: 6 3; animation: nm-dash 1.4s linear infinite; }
@keyframes nm-dash { to { stroke-dashoffset: -32; } }
.nm-svg-label { font-family: var(--nm-mono); font-size: 10px; fill: var(--nm-fg-2); }
.nm-svg-teal { fill: var(--nm-teal); font-weight: 600; }
.nm-svg-dim { fill: var(--nm-fg-3); font-family: var(--nm-mono); font-size: 10px; }

.nm-packet { fill: var(--nm-teal); filter: drop-shadow(0 0 6px var(--nm-teal)); }
.nm-packet-1 { offset-path: path("M194 80 C 240 80, 240 160, 280 160 C 320 160, 320 70, 420 70"); animation: nm-packet-run-fast 1.6s linear infinite; }
.nm-packet-2 { offset-path: path("M194 140 C 240 140, 240 160, 280 160 C 320 160, 320 160, 420 160"); animation: nm-packet-run-fast 1.6s linear infinite; }
.nm-packet-3 { offset-path: path("M194 200 C 240 200, 240 160, 280 160 C 320 160, 320 250, 420 250"); animation: nm-packet-run-fast 1.6s linear infinite; }
@keyframes nm-packet-run-fast {
  0%   { offset-distance: 0%;  opacity: 0; }
  10%  { opacity: 1; }
  90%  { opacity: 1; }
  100% { offset-distance: 100%; opacity: 0; }
}

.nm-prov-svg { fill: #03050a; stroke: rgba(255,255,255,0.08); transition: all 0.2s; }
.nm-prov-svg-active { stroke: var(--nm-teal); fill: rgba(20, 184, 166, 0.08); }

.nm-routing-target {
  margin-top: 14px;
  display: flex; flex-wrap: wrap; align-items: center; gap: 8px;
  font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-2);
}
.nm-rt-label { color: var(--nm-fg-3); letter-spacing: 0.06em; text-transform: uppercase; }
.nm-rt-val {
  color: var(--nm-teal); background: rgba(20, 184, 166, 0.08); border: 1px solid rgba(20, 184, 166, 0.15);
  padding: 3px 8px; border-radius: 4px; font-feature-settings: "tnum";
}

.nm-log-feed {
  margin-top: 14px; background: #03050a; border: 1px solid var(--nm-line);
  border-radius: 8px; padding: 10px 14px; font-family: var(--nm-mono); font-size: 11px;
  height: 84px; overflow: hidden; position: relative;
}
.nm-log-feed::before {
  content: "$ tail -f /var/log/nexusmind/routing.log";
  display: block; font-size: 10px; color: var(--nm-fg-3);
  letter-spacing: 0.06em; padding-bottom: 6px;
  border-bottom: 1px dashed var(--nm-line); margin-bottom: 6px;
}
.nm-log-line {
  display: flex; gap: 8px; align-items: baseline; line-height: 1.5;
  animation: nm-log-in 0.4s ease; white-space: nowrap; overflow: hidden;
}
@keyframes nm-log-in {
  from { opacity: 0; transform: translateY(8px); }
  to   { opacity: 1; transform: translateY(0); }
}
.nm-log-t { color: var(--nm-fg-3); }
.nm-log-method { color: var(--nm-violet); }
.nm-log-arr { color: var(--nm-fg-3); }
.nm-log-prov { color: var(--nm-teal); }
.nm-log-model { color: var(--nm-fg-1); }
.nm-log-latency { color: var(--nm-fg-2); margin-left: auto; }
.nm-log-ok { color: var(--nm-emerald); }

.nm-console-meta { margin-top: 18px; display: grid; grid-template-columns: repeat(3,1fr); gap: 12px; font-family: var(--nm-mono); font-size: 11px; }
.nm-meta-cell { background: #03050a; border: 1px solid var(--nm-line); border-radius: 8px; padding: 10px 12px; }
.nm-meta-k { color: var(--nm-fg-3); letter-spacing: .08em; text-transform: uppercase; font-size: 10px; }
.nm-meta-v { color: var(--nm-fg); font-size: 16px; margin-top: 4px; letter-spacing: -.01em; font-feature-settings: "tnum"; font-weight: 600; }
.nm-meta-u { color: var(--nm-fg-3); font-size: 11px; margin-left: 2px; }

/* ── Video Tutorial Player UI ── */
.nm-tut-grid { display: grid; grid-template-columns: 1.65fr 1fr; gap: 36px; margin-top: 48px; }
@media (max-width: 920px) { .nm-tut-grid { grid-template-columns: 1fr; gap: 28px; } }
.nm-video-wrap {
  position: relative; aspect-ratio: 16 / 9; width: 100%; border-radius: var(--nm-r);
  overflow: hidden; background: #060910; border: 1px solid var(--nm-line-2);
  box-shadow: 0 20px 45px rgba(0,0,0,0.6), inset 0 1px 0 rgba(255,255,255,0.05);
}
.nm-video-wrap iframe { position: absolute; inset: 0; width: 100%; height: 100%; border: 0; }
.nm-video-placeholder {
  position: absolute; inset: 0; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 18px;
  background: radial-gradient(600px 300px at 50% 60%, rgba(20, 184, 166, 0.06), transparent 60%),
              repeating-linear-gradient(45deg, #05080f 0px, #05080f 12px, #080c15 12px, #080c15 24px);
  cursor: pointer; transition: background 0.2s;
}
.nm-video-placeholder:hover {
  background: radial-gradient(700px 350px at 50% 60%, rgba(20, 184, 166, 0.12), transparent 60%),
              repeating-linear-gradient(45deg, #05080f 0px, #05080f 12px, #080c15 12px, #080c15 24px);
}
.nm-video-badge {
  position: absolute; top: 18px; left: 18px; display: inline-flex; align-items: center; gap: 8px;
  font-family: var(--nm-mono); font-size: 11px; letter-spacing: 0.14em; text-transform: uppercase;
  color: var(--nm-teal); background: rgba(3,7,18,0.85); border: 1px solid var(--nm-line-2);
  padding: 6px 12px; border-radius: 999px;
}
.nm-video-time {
  position: absolute; top: 18px; right: 18px; font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-2);
  background: rgba(3,7,18,0.85); border: 1px solid var(--nm-line-2); padding: 6px 12px; border-radius: 999px;
}
.nm-play-btn {
  width: 78px; height: 78px; border-radius: 50%; background: var(--nm-teal); color: #020617;
  display: grid; place-items: center; border: none; cursor: pointer;
  animation: nm-play-pulse 2s ease infinite; transition: transform 0.15s;
}
.nm-play-btn:hover { transform: scale(1.08); background: #2dd4bf; }
.nm-play-btn svg { margin-left: 6px; }
@keyframes nm-play-pulse {
  0%, 100% { box-shadow: 0 0 0 0 rgba(20, 184, 166, 0.45); }
  50%      { box-shadow: 0 0 0 16px rgba(20, 184, 166, 0); }
}
.nm-video-hint { font-family: var(--nm-mono); font-size: 12px; color: var(--nm-fg-2); letter-spacing: 0.08em; text-transform: uppercase; text-align: center; padding: 0 24px; }
.nm-video-hint code { color: var(--nm-teal); }
.nm-video-strip {
  position: absolute; bottom: 16px; left: 18px; right: 18px;
  display: flex; gap: 14px; align-items: center; justify-content: space-between;
  font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-3);
}
.nm-vs-left { display: inline-flex; align-items: center; gap: 8px; }
.nm-vs-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--nm-teal); display: inline-block; }
.nm-vs-progress { flex: 1; max-width: 200px; height: 3px; background: rgba(255,255,255,0.08); border-radius: 2px; overflow: hidden; position: relative; }
.nm-vs-progress::after { content: ""; display: block; width: 25%; height: 100%; background: var(--nm-teal); }

.nm-tut-h { font-family: var(--nm-mono); font-size: 11px; letter-spacing: 0.14em; text-transform: uppercase; color: var(--nm-fg-3); margin: 0 0 16px; }
.nm-tut-list { list-style: none; margin: 0; padding: 0; }
.nm-tut-list li {
  display: grid; grid-template-columns: 56px 1fr; align-items: baseline; gap: 14px;
  padding: 14px 0; border-top: 1px solid var(--nm-line); font-size: 15px; color: var(--nm-fg);
  transition: padding-left 0.15s, color 0.15s; cursor: pointer;
}
.nm-tut-list li:hover { padding-left: 8px; color: var(--nm-teal); }
.nm-tut-list li:last-child { border-bottom: 1px solid var(--nm-line); }
.nm-tut-time { font-family: var(--nm-mono); font-size: 12px; color: var(--nm-teal); letter-spacing: 0.04em; font-feature-settings: "tnum"; font-weight: 600; }

/* ── Action Banner ── */
.nm-action-banner {
  margin-top: 40px; padding: 28px; background: linear-gradient(90deg, rgba(20, 184, 166, 0.05) 0%, rgba(124, 58, 237, 0.05) 100%);
  border: 1px solid var(--nm-line-2); border-radius: var(--nm-r); display: flex; gap: 24px; align-items: center;
  justify-content: space-between; box-shadow: 0 15px 35px rgba(0,0,0,0.4);
}
.nm-tut-grid .nm-action-banner {
  margin-top: 24px;
}
.nm-ab-content { display: flex; gap: 20px; align-items: center; }
.nm-ab-icon { font-size: 32px; flex-shrink: 0; animation: bounce 3s infinite; }
.nm-ab-text h4 { font-family: var(--nm-title); font-size: 18px; font-weight: 700; margin: 0 0 4px; color: var(--nm-fg); }
.nm-ab-text p { font-size: 14px; color: var(--nm-fg-2); margin: 0; line-height: 1.5; }
@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}
@media (max-width: 768px) {
  .nm-action-banner { flex-direction: column; text-align: center; gap: 20px; padding: 24px; }
  .nm-ab-content { flex-direction: column; gap: 12px; }
  .nm-action-banner .nm-btn { width: 100%; justify-content: center; }
}

/* ── Steps ── */
.nm-steps { display: grid; grid-template-columns: repeat(3,1fr); gap: 20px; margin-top: 48px; }
@media (max-width: 920px) { .nm-steps { grid-template-columns: 1fr; } }
.nm-step { background: var(--nm-bg-1); border: 1px solid var(--nm-line); border-radius: var(--nm-r); padding: 28px; transition: border-color .2s; }
.nm-step:hover { border-color: var(--nm-line-2); }
.nm-step-num { font-family: var(--nm-mono); font-size: 12px; color: var(--nm-fg-3); letter-spacing: .1em; }
.nm-step-num b { color: var(--nm-teal); font-weight: 600; }
.nm-step h3 { font-family: var(--nm-title); font-size: 22px; font-weight: 700; letter-spacing: -.015em; margin: 14px 0 8px; }
.nm-step p { color: var(--nm-fg-1); font-size: 15px; line-height: 1.55; margin: 0 0 18px; }
.nm-codeblock { background: #03050a; border: 1px solid var(--nm-line); font-family: var(--nm-mono); font-size: 12px; padding: 12px 14px; border-radius: 8px; color: var(--nm-fg-1); overflow: hidden; }
.nm-c { color: var(--nm-fg-3); }
.nm-k { color: var(--nm-teal); }
.nm-v { color: var(--nm-violet); }
.nm-s { color: var(--nm-emerald); }

/* ── Matemática ── */
.nm-math-grid { display: grid; grid-template-columns: 1.1fr 1fr; gap: 48px; align-items: center; margin-top: 56px; }
@media (max-width: 1000px) { .nm-math-grid { grid-template-columns: 1fr; gap: 40px; } }
.nm-compare { background: var(--nm-bg-1); border: 1px solid var(--nm-line); border-radius: var(--nm-r); overflow: hidden; }
.nm-cmp-row { display: grid; grid-template-columns: 1.2fr 1fr 1fr; border-bottom: 1px solid var(--nm-line); font-family: var(--nm-mono); font-size: 13px; }
.nm-cmp-row:last-child { border-bottom: none; }
.nm-cmp-header { background: #05080f; }
.nm-cmp-header .nm-cmp-cell { font-size: 10px; letter-spacing: .1em; text-transform: uppercase; color: var(--nm-fg-3); padding: 14px 18px; }
.nm-cmp-header .nm-cmp-ours { color: var(--nm-teal); }
.nm-cmp-cell { padding: 16px 18px; border-right: 1px solid var(--nm-line); color: var(--nm-fg-1); }
.nm-cmp-cell:last-child { border-right: none; }
.nm-cmp-label { color: var(--nm-fg-2); }
.nm-cmp-ours { color: var(--nm-fg); font-weight: 500; background: rgba(20, 184, 166, 0.02); }
.nm-cmp-total { background: #05080f; }
.nm-cmp-total .nm-cmp-cell { font-size: 15px; color: var(--nm-fg); padding: 18px; font-weight: 700; font-family: var(--nm-title); }
.nm-cmp-total .nm-cmp-ours { color: var(--nm-teal); }

.nm-math-points { display: flex; flex-direction: column; gap: 24px; }
.nm-math-pt { border-left: 2px solid var(--nm-line-2); padding: 4px 0 4px 20px; cursor: default; transition: border-left-color 0.2s; }
.nm-math-pt-active { border-left-color: var(--nm-teal); }
.nm-pt-tag { font-family: var(--nm-mono); font-size: 11px; letter-spacing: .1em; color: var(--nm-violet); text-transform: uppercase; }
.nm-math-pt h4 { font-family: var(--nm-title); font-size: 18px; font-weight: 700; letter-spacing: -.01em; margin: 6px 0 8px; }
.nm-math-pt p { color: var(--nm-fg-1); font-size: 14.5px; line-height: 1.55; margin: 0; }

/* ── Integrations ── */
.nm-int-grid { display: grid; grid-template-columns: repeat(4,1fr); gap: 14px; margin-top: 48px; }
@media (max-width: 1040px) { .nm-int-grid { grid-template-columns: repeat(2,1fr); } }
@media (max-width: 540px)  { .nm-int-grid { grid-template-columns: 1fr; } }
.nm-int-card {
  background: var(--nm-bg-1); border: 1px solid var(--nm-line); border-radius: var(--nm-r);
  padding: 20px; position: relative; transition: border-color 0.2s, transform 0.2s, background 0.2s;
}
.nm-int-card:hover {
  border-color: var(--nm-teal);
  background: linear-gradient(180deg, rgba(20, 184, 166, 0.05) 0%, #0b0f19 100%);
  transform: translateY(-2px);
}
.nm-int-row { display: flex; align-items: center; gap: 12px; margin-bottom: 14px; }
.nm-int-icon { width: 36px; height: 36px; background: #03050a; border: 1px solid var(--nm-line-2); border-radius: 8px; display: grid; place-items: center; font-family: var(--nm-mono); font-size: 16px; font-weight: 700; color: var(--nm-teal); }
.nm-int-name { font-family: var(--nm-title); font-size: 16px; font-weight: 700; letter-spacing: -0.01em; }
.nm-int-where { font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-3); letter-spacing: 0.04em; text-transform: uppercase; margin-bottom: 10px; }
.nm-int-env { background: #03050a; border: 1px solid var(--nm-line); font-family: var(--nm-mono); font-size: 11px; padding: 8px 10px; border-radius: 6px; color: var(--nm-fg-2); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.nm-int-env b { color: var(--nm-teal); font-weight: 600; }
.nm-int-url { color: var(--nm-violet); }
.nm-int-tested {
  position: absolute; top: 18px; right: 18px; font-family: var(--nm-mono); font-size: 9px; color: var(--nm-emerald);
  display: inline-flex; align-items: center; gap: 4px; letter-spacing: 0.1em; text-transform: uppercase;
}
.nm-int-d { width: 5px; height: 5px; border-radius: 50%; background: var(--nm-emerald); box-shadow: 0 0 6px var(--nm-emerald); display: inline-block; }
.nm-int-more {
  margin-top: 32px; padding: 20px 24px; border: 1px dashed var(--nm-line-2); border-radius: var(--nm-r);
  font-family: var(--nm-mono); font-size: 12px; color: var(--nm-fg-2); letter-spacing: 0.04em; line-height: 1.6;
}
.nm-int-more b { color: var(--nm-teal); font-weight: 600; }

/* ── Providers ── */
.nm-providers { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; margin-top: 48px; }
@media (max-width: 1040px) { .nm-providers { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 640px) { .nm-providers { grid-template-columns: 1fr; } }
.nm-prov-card { background: var(--nm-bg-1); border: 1px solid var(--nm-line); border-radius: var(--nm-r); padding: 22px; position: relative; transition: all .2s; }
.nm-prov-card:hover { border-color: var(--nm-line-2); transform: translateY(-2px); }
.nm-prov-status { position: absolute; top: 18px; right: 18px; font-family: var(--nm-mono); font-size: 10px; color: var(--nm-emerald); display: inline-flex; align-items: center; gap: 6px; }
.nm-prov-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--nm-emerald); box-shadow: 0 0 8px var(--nm-emerald); display: inline-block; }
.nm-prov-ico { width: 40px; height: 40px; border-radius: 10px; display: grid; place-items: center; font-family: var(--nm-mono); font-size: 16px; font-weight: 700; background: #03050a; border: 1px solid; margin-bottom: 14px; }
.nm-prov-card h4 { font-family: var(--nm-title); font-size: 16px; font-weight: 700; margin: 0 0 6px; letter-spacing: -.01em; }
.nm-prov-models { font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-3); letter-spacing: .04em; line-height: 1.6; }
.nm-prov-models b { color: var(--nm-fg); font-weight: 500; }
.nm-prov-pills { margin-top: 10px; display: flex; flex-wrap: wrap; gap: 6px; }
.nm-prov-pills span { display: inline-block; padding: 3px 8px; border-radius: 999px; background: rgba(255,255,255,0.04); border: 1px solid rgba(255,255,255,0.06); color: var(--nm-fg-1); font-size: 10.5px; transition: all 0.2s ease; }
.nm-prov-card:hover .nm-prov-pills span { border-color: rgba(20, 184, 166, 0.15); background: rgba(20, 184, 166, 0.02); }
.nm-providers-footer { margin-top: 32px; padding: 16px 20px; border: 1px dashed var(--nm-line-2); border-radius: var(--nm-r); font-family: var(--nm-mono); font-size: 12px; color: var(--nm-fg-2); line-height: 1.6; }
.nm-providers-footer b { color: var(--nm-teal); font-weight: 600; }
.animate-ping-slow { animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite; }

/* ── Pricing ── */
.nm-pricing-grid { display: grid; grid-template-columns: repeat(3,1fr); gap: 20px; margin-top: 48px; align-items: stretch; }
@media (max-width: 960px) { .nm-pricing-grid { grid-template-columns: 1fr; } }
.nm-plan { background: var(--nm-bg-1); border: 1px solid var(--nm-line); border-radius: var(--nm-r); padding: 32px; position: relative; display: flex; flex-direction: column; }
.nm-plan-feat {
  border-color: var(--nm-teal);
  background: linear-gradient(180deg, rgba(20, 184, 166, 0.04) 0%, #0b0f19 100%);
  box-shadow: 0 20px 45px rgba(0,0,0,0.6);
}
.nm-plan-badge { position: absolute; top: -10px; left: 32px; font-family: var(--nm-mono); font-size: 10px; letter-spacing: .14em; text-transform: uppercase; background: var(--nm-teal); color: #020617; padding: 4px 10px; border-radius: 4px; font-weight: 700; }
.nm-plan-type { font-family: var(--nm-mono); font-size: 11px; letter-spacing: .14em; text-transform: uppercase; color: var(--nm-fg-3); }
.nm-plan h3 { font-family: var(--nm-title); font-size: 24px; font-weight: 700; letter-spacing: -.015em; margin: 0 0 8px; }
.nm-plan-desc { font-size: 14px; color: var(--nm-fg-2); line-height: 1.5; margin: 0 0 22px; min-height: 42px; }
.nm-plan-price { display: flex; align-items: baseline; gap: 4px; margin-bottom: 6px; }
.nm-price-cur { font-family: var(--nm-mono); font-size: 14px; color: var(--nm-fg-2); }
.nm-price-amt { font-family: var(--nm-title); font-size: 56px; font-weight: 800; letter-spacing: -.04em; line-height: 1; }
.nm-price-per { font-family: var(--nm-mono); font-size: 13px; color: var(--nm-fg-3); margin-left: 4px; }
.nm-plan-old { font-family: var(--nm-mono); font-size: 12px; color: var(--nm-fg-3); text-decoration: line-through; margin-bottom: 22px; }
.nm-plan-old-spacer { text-decoration: none; visibility: hidden; }
.nm-plan-features { list-style: none; margin: 0 0 24px; padding: 0; flex: 1; }
.nm-plan-features li { font-size: 14px; color: var(--nm-fg-1); padding: 8px 0 8px 24px; position: relative; border-bottom: 1px solid var(--nm-line); }
.nm-plan-features li:last-child { border-bottom: none; }
.nm-plan-features li::before { content: "→"; position: absolute; left: 0; top: 8px; font-family: var(--nm-mono); color: var(--nm-teal); font-size: 14px; font-weight: 700; }
.nm-plan-features li.nm-feat-muted { color: var(--nm-fg-3); }
.nm-plan-features li.nm-feat-muted::before { content: "·"; color: var(--nm-fg-3); }

/* ── Testimonials ── */
.nm-proof-grid { display: grid; grid-template-columns: repeat(3,1fr); gap: 20px; margin-top: 48px; }
@media (max-width: 1000px) { .nm-proof-grid { grid-template-columns: 1fr; } }
.nm-tcard { background: var(--nm-bg-1); border: 1px solid var(--nm-line); border-radius: var(--nm-r); padding: 28px; display: flex; flex-direction: column; }
.nm-tcard-meta { font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-3); margin-bottom: 16px; letter-spacing: .08em; text-transform: uppercase; display: flex; gap: 8px; }
.nm-t-sep { color: var(--nm-line-2); }
.nm-t-ec { color: var(--nm-emerald); }
.nm-tcard-quote { font-size: 16px; line-height: 1.55; letter-spacing: -.01em; color: var(--nm-fg); margin: 0 0 24px; text-wrap: pretty; flex: 1; }
.nm-tcard-who { display: flex; align-items: center; gap: 12px; }
.nm-tcard-av {
  width: 36px; height: 36px; border-radius: 50%;
  background: linear-gradient(135deg, var(--nm-teal), var(--nm-violet));
  display: grid; place-items: center; font-family: var(--nm-mono); font-size: 13px; font-weight: 700; color: #020617; flex-shrink: 0;
}
.nm-tcard-who b { display: block; font-size: 14px; font-weight: 600; font-family: var(--nm-title); }
.nm-tcard-who span { display: block; font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-3); }

/* ── FAQ ── */
.nm-faq-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 0 36px; margin-top: 32px; }
@media (max-width: 880px) { .nm-faq-grid { grid-template-columns: 1fr; } }
.nm-faq-item { border-top: 1px solid var(--nm-line); }
.nm-faq-item button {
  display: flex; align-items: center; justify-content: space-between; gap: 18px;
  width: 100%; padding: 20px 0; background: transparent; border: none; cursor: pointer; text-align: left;
  color: var(--nm-fg); font-size: 15.5px; font-weight: 600; font-family: inherit;
}
.nm-faq-item button:hover { color: var(--nm-teal); }
.nm-faq-num { font-family: var(--nm-mono); font-size: 12px; color: var(--nm-fg-3); margin-right: 14px; }
.nm-faq-plus { color: var(--nm-teal); font-family: var(--nm-mono); font-size: 18px; transition: transform .2s; flex-shrink: 0; }
.nm-faq-open .nm-faq-plus { transform: rotate(45deg); }
.nm-faq-answer { max-height: 0; overflow: hidden; transition: max-height .25s ease, padding .25s ease; }
.nm-faq-open .nm-faq-answer { max-height: 400px; padding: 0 0 24px 30px; }
.nm-faq-answer p { color: var(--nm-fg-2); font-size: 14.5px; line-height: 1.6; margin: 0; max-width: 56ch; }

/* ── Final CTA ── */
.nm-final-cta {
  background: var(--nm-bg-1); border: 1px solid var(--nm-line); border-radius: var(--nm-r);
  padding: 56px; text-align: center; position: relative; overflow: hidden;
}
.nm-final-cta::before {
  content: ""; position: absolute; inset: 0; pointer-events: none;
  background: radial-gradient(600px 300px at 50% 0%, rgba(20, 184, 166, 0.1), transparent 70%);
}
.nm-final-cta .nm-s-label { display: inline-flex; }
.nm-final-cta .nm-s-title { font-size: clamp(32px,4vw,48px); position: relative; }
.nm-cta-p { font-size: 16px; color: var(--nm-fg-2); margin: 0 auto 28px; max-width: 520px; position: relative; line-height: 1.55; }
.nm-cta-btns { display: inline-flex; gap: 12px; position: relative; flex-wrap: wrap; justify-content: center; }
@media (max-width: 720px) { .nm-final-cta { padding: 40px 24px; } }

/* ── Footer ── */
.nm-footer { border-top: 1px solid var(--nm-line); margin-top: 80px; padding: 64px 0 28px; position: relative; z-index: 1; }
.nm-footer-grid { display: grid; grid-template-columns: 2fr 1fr 1fr 1fr; gap: 36px; }
@media (max-width: 880px) { .nm-footer-grid { grid-template-columns: 1fr 1fr; } }
.nm-footer-col h5 { font-family: var(--nm-mono); font-size: 11px; letter-spacing: .14em; text-transform: uppercase; color: var(--nm-fg-3); margin: 0 0 16px; }
.nm-footer-col a { display: block; font-size: 14.5px; color: var(--nm-fg-2); padding: 6px 0; text-decoration: none; transition: color 0.15s; }
.nm-footer-col a:hover { color: var(--nm-teal); }
.nm-footer-desc { font-size: 14px; color: var(--nm-fg-3); line-height: 1.55; max-width: 360px; margin: 12px 0 0; }
.nm-footer-bot {
  display: flex; justify-content: space-between; align-items: center;
  margin-top: 48px; padding-top: 24px; border-top: 1px solid var(--nm-line);
  font-family: var(--nm-mono); font-size: 11px; color: var(--nm-fg-3); letter-spacing: .04em;
}
@media (max-width: 720px) { .nm-footer-bot { flex-direction: column; gap: 12px; text-align: center; } }

/* ── Floating CTA ── */
.nm-float-cta {
  position: fixed; left: 50%; bottom: 24px;
  transform: translateX(-50%) translateY(120%);
  z-index: 80; transition: transform 0.35s cubic-bezier(0.5, 1.5, 0.3, 1);
  pointer-events: none;
}
.nm-float-visible { transform: translateX(-50%) translateY(0); pointer-events: auto; }
.nm-float-inner {
  display: flex; align-items: center; gap: 16px; padding: 6px 6px 6px 16px;
  background: rgba(11, 15, 25, 0.9); backdrop-filter: blur(16px);
  border: 1px solid var(--nm-line-2); border-radius: 999px;
  box-shadow: 0 20px 50px rgba(0,0,0,0.6);
  font-family: var(--nm-mono); font-size: 12px; color: var(--nm-fg-1);
}
.nm-float-prompt { color: var(--nm-teal); font-weight: 700; }
.nm-float-sep { color: var(--nm-line-2); margin: 0 -8px; }
.nm-float-link {
  color: var(--nm-fg-2); letter-spacing: 0.04em; text-transform: uppercase; font-size: 11px;
  transition: color 0.15s; text-decoration: none;
}
.nm-float-link:hover { color: var(--nm-teal); }
.nm-float-btn {
  height: 36px; padding: 0 16px; background: var(--nm-teal); color: #020617;
  border: none; border-radius: 999px; font-family: var(--nm-mono); font-size: 12px; font-weight: 600;
  letter-spacing: 0.04em; text-transform: uppercase; display: inline-flex; align-items: center; gap: 6px;
  cursor: pointer; text-decoration: none; transition: background 0.15s;
}
.nm-float-btn:hover { background: #2dd4bf; }
.nm-float-close {
  width: 32px; height: 32px; display: grid; place-items: center;
  background: transparent; border: none; color: var(--nm-fg-3); cursor: pointer;
  border-radius: 50%; font-family: var(--nm-sans); font-size: 16px;
}
.nm-float-close:hover { color: var(--nm-fg); background: rgba(255,255,255,0.05); }
@media (max-width: 720px) {
  .nm-float-link, .nm-float-sep { display: none; }
}

/* ── Scroll Reveal ── */
.nm-reveal {
  opacity: 0; transform: translateY(20px);
  transition: opacity 0.6s ease, transform 0.6s ease;
}
.nm-reveal-in { opacity: 1; transform: translateY(0); }
@media (prefers-reduced-motion: reduce) {
  .nm-reveal, .nm-reveal-in { opacity: 1; transform: none; transition: none; }
}
</style>
