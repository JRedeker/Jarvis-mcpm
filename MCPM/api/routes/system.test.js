/**
 * System Routes Tests
 */

const request = require('supertest');
const { createApp } = require('../server');

describe('System Routes', () => {
    let app;

    beforeAll(() => {
        app = createApp();
    });

    describe('GET /api/v1/usage', () => {
        it('should return usage statistics', async () => {
            const res = await request(app)
                .get('/api/v1/usage')
                .expect(200);

            expect(res.body.success).toBe(true);
            expect(res.body.data).toBeDefined();
            expect(res.body.data.servers).toBeDefined();
            expect(res.body.data.servers.total).toBeGreaterThanOrEqual(0);
            expect(res.body.data.profiles).toBeDefined();
            expect(res.body.data.clients).toBeDefined();
            expect(res.body.data.configDir).toBeDefined();
        });

        it('should include server breakdown', async () => {
            const res = await request(app)
                .get('/api/v1/usage')
                .expect(200);

            expect(res.body.data.servers.byCategory).toBeDefined();
            expect(typeof res.body.data.servers.byCategory).toBe('object');
        });
    });

    describe('POST /api/v1/migrate', () => {
        it('should perform migration', async () => {
            const res = await request(app)
                .post('/api/v1/migrate')
                .expect(200);

            expect(res.body.success).toBe(true);
            expect(res.body.data).toBeDefined();
            expect(Array.isArray(res.body.data.migrations)).toBe(true);
            expect(res.body.data.configDir).toBeDefined();
            expect(res.body.data.message).toBe('Migration complete');
        });
    });

    describe('GET /api/v1/search', () => {
        it('should search servers', async () => {
            const res = await request(app)
                .get('/api/v1/search?q=test')
                .expect(200);

            expect(res.body.success).toBe(true);
            expect(res.body.data).toBeDefined();
            expect(res.body.data.query).toBe('test');
            expect(res.body.data.count).toBeGreaterThanOrEqual(0);
            expect(Array.isArray(res.body.data.results)).toBe(true);
        });

        it('should return 400 for missing query', async () => {
            const res = await request(app)
                .get('/api/v1/search')
                .expect(400);

            expect(res.body.success).toBe(false);
            expect(res.body.error.code).toBe('VALIDATION_ERROR');
        });
    });

    describe('404 Handler', () => {
        it('should return 404 for unknown endpoints', async () => {
            const res = await request(app)
                .get('/api/v1/unknown-endpoint')
                .expect(404);

            expect(res.body.success).toBe(false);
            expect(res.body.error.code).toBe('NOT_FOUND');
        });

        it('should return 404 for wrong HTTP method', async () => {
            const res = await request(app)
                .patch('/api/v1/health')
                .expect(404);

            expect(res.body.success).toBe(false);
        });
    });
});
