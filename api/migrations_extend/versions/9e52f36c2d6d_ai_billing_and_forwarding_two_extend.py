"""ai billing and forwarding two extend

Revision ID: 004_ai_billing_forwarding
Revises: 003_tenant_model_sync_extend
Create Date: 2024-07-23 19:27:09.150306

"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.engine.reflection import Inspector

from models import db, types

# revision identifiers, used by Alembic.
revision = '004_ai_billing_forwarding'
down_revision = '003_tenant_model_sync_extend'
branch_labels = None
depends_on = None


def upgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    conn = op.get_bind()
    inspector = Inspector.from_engine(conn)
    tables = inspector.get_table_names()
    
    if 'forwarding_address_extend' not in tables:
        op.create_table('forwarding_address_extend',
        sa.Column('id', types.StringUUID(), server_default=sa.text('uuid_generate_v4()'), nullable=False),
        sa.Column('forwarding_id', types.StringUUID(), nullable=False),
        sa.Column('path', sa.String(length=255), nullable=False),
        sa.Column('models', sa.String(length=255), nullable=False),
        sa.Column('description', sa.Text(), server_default=sa.text("''::character varying"), nullable=False),
        sa.Column('content_type', sa.Integer(), server_default=sa.text('0'), nullable=False),
        sa.Column('billing', sa.Text(), server_default=sa.text("'[]'"), nullable=False),
        sa.Column('status', sa.Boolean(), server_default=sa.text('true'), nullable=True),
        sa.PrimaryKeyConstraint('id', name='forwarding_address_pkey')
        )
        with op.batch_alter_table('forwarding_address_extend', schema=None) as batch_op:
            batch_op.create_index('idx_forwarding_address_id', ['forwarding_id'], unique=False)
            batch_op.create_index('idx_forwarding_address_path', ['path'], unique=False)
            batch_op.create_index('idx_forwarding_address_status', ['status'], unique=False)

    if 'forwarding_extend' not in tables:
        op.create_table('forwarding_extend',
        sa.Column('id', types.StringUUID(), server_default=sa.text('uuid_generate_v4()'), nullable=False),
        sa.Column('path', sa.String(length=255), nullable=False),
        sa.Column('address', sa.String(length=255), nullable=False),
        sa.Column('header', sa.Text(), server_default=sa.text("'[]'"), nullable=False),
        sa.Column('description', sa.Text(), server_default=sa.text("''::character varying"), nullable=False),
        sa.PrimaryKeyConstraint('id', name='forwarding_extend_pkey')
        )
        with op.batch_alter_table('forwarding_extend', schema=None) as batch_op:
            batch_op.create_index('idx_forwarding_path', ['path'], unique=False)
    # ### end Alembic commands ###


def downgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    with op.batch_alter_table('forwarding_extend', schema=None) as batch_op:
        batch_op.drop_index('idx_forwarding_path')

    op.drop_table('forwarding_extend')
    with op.batch_alter_table('forwarding_address_extend', schema=None) as batch_op:
        batch_op.drop_index('idx_forwarding_address_path')
        batch_op.drop_index('idx_forwarding_address_id')
        batch_op.drop_index('idx_forwarding_address_status')

    op.drop_table('forwarding_address_extend')
    # ### end Alembic commands ###
