/**
 * Servers Routes Tests
 */

const request = require('supertest');
const { createApp } = require('../server');

describe('Servers Routes', () => {
    let app;

    beforeAll(() => {
        app = createApp();
    });

    describe('GET /api/v1/servers', () => {
        it('should return list of servers', async () => {
            const res = await request(app)
                .get('/api/v1/servers')
                .expect(200);

            expect(res.body.success).toBe(true);
            expect(res.body.data).toBeDefined();
            expect(res.body.data.count).toBeGreaterThanOrEqual(0);
            expect(Array.isArray(res.body.data.servers)).toBe(true);
        });

        it('should include server details', async () => {
            const res = await request(app)
                .get('/api/v1/servers')
                .expect(200);

            if (res.body.data.servers.length > 0) {
                const server = res.body.data.servers[0];
                expect(server.name).toBeDefined();
                expect(server.group).toBeDefined();
                expect(typeof server.installed).toBe('boolean');
            }
        });
    });

    describe('GET /api/v1/servers/:name', () => {
        it('should return 404 for non-existent server', async () => {
            const res = await request(app)
                .get('/api/v1/servers/nonexistent-server-xyz')
                .expect(404);

            expect(res.body.success).toBe(false);
            expect(res.body.error.code).toBe('SERVER_NOT_FOUND');
        });
    });

    describe('GET /api/v1/servers/search', () => {
        it('should search servers', async () => {
            const res = await request(app)
                .get('/api/v1/servers/search?q=test')
                .expect(200);

            expect(res.body.success).toBe(true);
            expect(res.body.data).toBeDefined();
            expect(res.body.data.query).toBe('test');
            expect(res.body.data.count).toBeGreaterThanOrEqual(0);
            expect(Array.isArray(res.body.data.results)).toBe(true);
        });

        it('should return 400 for missing query', async () => {
            const res = await request(app)
                .get('/api/v1/servers/search')
                .expect(400);

            expect(res.body.success).toBe(false);
            expect(res.body.error.code).toBe('VALIDATION_ERROR');
        });

        it('should return 400 for empty query', async () => {
            const res = await request(app)
                .get('/api/v1/servers/search?q=')
                .expect(400);

            expect(res.body.success).toBe(false);
        });
    });

    describe('POST /api/v1/servers', () => {
        it('should return 400 for missing name', async () => {
            const res = await request(app)
                .post('/api/v1/servers')
                .send({ type: 'stdio', command: 'node' })
                .expect(400);

            expect(res.body.success).toBe(false);
            expect(res.body.error.code).toBe('VALIDATION_ERROR');
            expect(res.body.error.message).toContain('name');
        });

        it('should return 400 for missing type', async () => {
            const res = await request(app)
                .post('/api/v1/servers')
                .send({ name: 'test-server' })
                .expect(400);

            expect(res.body.success).toBe(false);
            expect(res.body.error.code).toBe('VALIDATION_ERROR');
        });

        it('should return 400 for invalid type', async () => {
            const res = await request(app)
                .post('/api/v1/servers')
                .send({ name: 'test', type: 'invalid' })
                .expect(400);

            expect(res.body.success).toBe(false);
        });

        it('should return 400 for SSE type (deprecated)', async () => {
            const res = await request(app)
                .post('/api/v1/servers')
                .send({ name: 'test', type: 'sse', url: 'http://localhost' })
                .expect(400);

            expect(res.body.success).toBe(false);
            expect(res.body.error.message).toContain('SSE');
        });

        it('should return 400 for stdio without command', async () => {
            const res = await request(app)
                .post('/api/v1/servers')
                .send({ name: 'test', type: 'stdio' })
                .expect(400);

            expect(res.body.success).toBe(false);
            expect(res.body.error.message).toContain('Command is required');
        });

        it('should return 400 for http without url', async () => {
            const res = await request(app)
                .post('/api/v1/servers')
                .send({ name: 'test', type: 'streamable-http' })
                .expect(400);

            expect(res.body.success).toBe(false);
            expect(res.body.error.message).toContain('URL is required');
        });
    });

    describe('PUT /api/v1/servers/:name', () => {
        it('should return 404 for non-custom server', async () => {
            const res = await request(app)
                .put('/api/v1/servers/nonexistent-xyz')
                .send({ command: 'node' })
                .expect(404);

            expect(res.body.success).toBe(false);
            expect(res.body.error.code).toBe('SERVER_NOT_FOUND');
        });
    });

    describe('DELETE /api/v1/servers/:name', () => {
        it('should return 404 for non-installed server', async () => {
            const res = await request(app)
                .delete('/api/v1/servers/nonexistent-server-xyz')
                .expect(404);

            expect(res.body.success).toBe(false);
        });
    });
});
