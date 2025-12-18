/**
 * Health Routes Tests
 */

const request = require('supertest');
const { createApp } = require('../server');

describe('Health Routes', () => {
    let app;

    beforeAll(() => {
        app = createApp();
    });

    describe('GET /api/v1/health', () => {
        it('should return health status', async () => {
            const res = await request(app)
                .get('/api/v1/health')
                .expect(200);

            expect(res.body.success).toBe(true);
            expect(res.body.data).toBeDefined();
            expect(res.body.data.status).toMatch(/^(healthy|degraded)$/);
            expect(res.body.data.timestamp).toBeDefined();
            expect(res.body.data.version).toBe('1.0.0');
            expect(res.body.data.checks).toBeDefined();
        });

        it('should include node check', async () => {
            const res = await request(app)
                .get('/api/v1/health')
                .expect(200);

            expect(res.body.data.checks.node).toBeDefined();
            expect(res.body.data.checks.node.status).toBe('ok');
            expect(res.body.data.checks.node.version).toMatch(/^v\d+/);
        });

        it('should include registry check', async () => {
            const res = await request(app)
                .get('/api/v1/health')
                .expect(200);

            expect(res.body.data.checks.registry).toBeDefined();
            expect(res.body.data.checks.registry.path).toBeDefined();
        });

        it('should include docker check', async () => {
            const res = await request(app)
                .get('/api/v1/health')
                .expect(200);

            expect(res.body.data.checks.docker).toBeDefined();
            expect(res.body.data.checks.docker.status).toMatch(/^(detected|not_found)$/);
        });
    });

    describe('GET /', () => {
        it('should redirect to health endpoint', async () => {
            const res = await request(app)
                .get('/')
                .expect(302);

            expect(res.headers.location).toBe('/api/v1/health');
        });
    });
});
